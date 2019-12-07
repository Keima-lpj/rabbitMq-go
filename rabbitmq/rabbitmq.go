package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

//rabbitMq的默认连接格式为：amqp://username:password@url/vhost
const MQ_URL = "amqp://immoc:immoc@127.0.0.1:5672/immoc"

//创建一个rabbitMQ的结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel

	//队列名称
	QueueName string
	//交换机名称
	Exchange string
	//路由key
	RouteKey string
}

//释放rabbitMq的资源
func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

//RabbitMq的错误处理函数
func (r *RabbitMQ) ErrorHandle(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//新建一个rabbitMq的结构体
func NewRabbitMQ(queueName, exchange, routeKey string) *RabbitMQ {
	rabbitMq := &RabbitMQ{QueueName: queueName, Exchange: exchange, RouteKey: routeKey}
	var err error
	rabbitMq.conn, err = amqp.Dial(MQ_URL)
	rabbitMq.ErrorHandle(err, "rabbitMq连接失败")
	rabbitMq.channel, err = rabbitMq.conn.Channel()
	rabbitMq.ErrorHandle(err, "rabbitMq获取channel失败")
	return rabbitMq
}

//申请队列
func (r *RabbitMQ) QueueDeclare() (amqp.Queue, error) {
	//申请队列，如果队列不存在则创建，存在则跳过
	q, err := r.channel.QueueDeclare(
		r.QueueName, //队列名
		//队列数据是否持久化
		false,
		//是否自动删除
		false,
		//是否具有排他性（除创建者外他人不可见）
		false,
		//是否阻塞（发送消息之后是否等待服务器响应）
		false,
		//其他额外的参数
		nil,
	)
	return q, err
}

//申请一个新的交换机
func (r *RabbitMQ) ExchangeDeclare(kind string) error {
	return r.channel.ExchangeDeclare(
		r.Exchange,
		//订阅发布的模式默认为fanout, 路由模式下为direct
		kind,
		false,
		false,
		false,
		false,
		nil,
	)
}

//新建一个简单模式（p->mq->c）的rabbitMq连接。 这里的exchange不传代表使用默认的交换机
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

//简单模式下的消息生产发送
func (r *RabbitMQ) SimplePublish(message string) {
	//1、申请队列
	_, err := r.QueueDeclare()
	if err != nil {
		fmt.Println("Publish申请队列失败:", err)
		return
	}

	//2、发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true，根据路由规则找不到对应的队列，则会将数据返还给发送者
		false,
		//如果为true，对应的队列没有消费者的话，则会将数据返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text:plain",
			Body:        []byte(message),
		},
	)

	if err != nil {
		fmt.Println("发送消息失败:", err)
		return
	}
}

//简单模式下的接收消息
func (r *RabbitMQ) SimpleConsume() {
	//1、申请队列
	_, err := r.QueueDeclare()
	if err != nil {
		fmt.Println("Consume申请队列失败:", err)
		return
	}

	//2、接收消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		//用以区分多个消费者
		"",
		//是否自动应答（消费完成后要给队列应答消息，否则队列会认定为消费失败，重新推送）
		true,
		//是否具有排他性
		false,
		//如果设置为true，则代表不能将同一个connection中发送的消息用这个connection来消费
		false,
		//是否阻塞，这里为阻塞
		false,
		//其他参数
		nil,
	)
	if err != nil {
		fmt.Println("Consume获取消费chan失败:", err)
		return
	}

	wait := make(chan bool)
	//开始处理消费数据
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			ProcessMessage(d.Body)
		}
	}()

	log.Printf("Waiting for message, To exit press CTRL + c")
	//这里暂时先阻塞
	<-wait
}

//这里是获取到消息后的处理逻辑
func ProcessMessage(message []byte) {
	fmt.Println(string(message))
}

//获取一个订阅模式下的RabbitMq实例
func NewRabbitMQPubSub(exchange string) *RabbitMQ {
	return NewRabbitMQ("", exchange, "")
}

//订阅模式下发布消息
func (r *RabbitMQ) PubSubPublish(message string) {
	//1、申请交换机
	err := r.ExchangeDeclare("fanout")
	if err != nil {
		fmt.Sprintf("申请交换机失败：%s\n", err)
		return
	}

	//2、发布消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		//如果为true，根据路由规则找不到对应的队列，则会将数据返还给发送者
		false,
		//如果为true，对应的队列没有消费者的话，则会将数据返回给发送者
		false,
		amqp.Publishing{
			ContentType: "text:plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		fmt.Println("发送消息失败:", err)
		return
	}
}

//订阅模式下接收消息
func (r *RabbitMQ) PubSubConsumer() {
	//1、申请交换机
	err := r.ExchangeDeclare("fanout")
	if err != nil {
		fmt.Sprintf("申请交换机失败：%s\n", err)
		return
	}
	//2、申请随机的队列（如果有指定的可以用指定的队列）
	q, err := r.QueueDeclare()
	if err != nil {
		fmt.Sprintf("申请随机队列失败：%s\n", err)
		return
	}
	//3、绑定交换机和队列 队列名称，路由key，exchange名称
	err = r.channel.QueueBind(q.Name, "", r.Exchange, false, nil)
	if err != nil {
		fmt.Sprintf("绑定队列到交换机失败：%s\n", err)
		return
	}
	//4、开始消费队列
	msgs, err := r.channel.Consume(
		q.Name,
		//用以区分多个消费者
		"",
		//是否自动应答（消费完成后要给队列应答消息，否则队列会认定为消费失败，重新推送）
		true,
		//是否具有排他性
		false,
		//如果设置为true，则代表不能将同一个connection中发送的消息用这个connection来消费
		false,
		//是否阻塞，这里为阻塞
		false,
		//其他参数
		nil,
	)
	if err != nil {
		fmt.Println("Consume获取消费chan失败:", err)
		return
	}

	wait := make(chan bool)
	//开始处理消费数据
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			ProcessMessage(d.Body)
		}
	}()

	log.Printf("Waiting for message, To exit press CTRL + c")
	//这里暂时先阻塞
	<-wait
}

//获取一个路由模式下的RabbitMq实例
func NewRabbitMQRoute(exchange, routeKey string) *RabbitMQ {
	return NewRabbitMQ("", exchange, routeKey)
}

//路由模式下发布消息
func (r *RabbitMQ) RoutePublish(message string) {
	//1、申请一个交换机
	err := r.ExchangeDeclare("direct")
	if err != nil {
		fmt.Sprintf("新建交换机失败：%s", err)
		return
	}

	//2、发布消息，到具体的路由
	err = r.channel.Publish(
		r.Exchange,
		r.RouteKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

//路由模式下消费消息
func (r *RabbitMQ) RouteConsumer() {
	//1、申请一个交换机
	err := r.ExchangeDeclare("direct")
	if err != nil {
		fmt.Sprintf("申请交换机失败：%s\n", err)
		return
	}
	//2、申请一个队列
	q, err := r.QueueDeclare()
	if err != nil {
		fmt.Sprintf("新建队列失败：%s\n", err)
		return
	}
	//3、绑定队列到交换机，注意这里要指定绑定的路由key
	err = r.channel.QueueBind(q.Name, r.RouteKey, r.Exchange, false, nil)
	if err != nil {
		fmt.Sprintf("绑定队列到交换机失败：%s\n", err)
	}
	//4、消费消息
	msgs, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)

	wait := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			ProcessMessage(d.Body)
		}
	}()

	log.Printf("Waiting for message, To exit press CTRL + c")
	//这里暂时先阻塞
	<-wait
}
