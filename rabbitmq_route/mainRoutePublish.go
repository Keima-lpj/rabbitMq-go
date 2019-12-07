package main

import (
	"fmt"
	"immoc-rabbitmq/rabbitmq"
	"strconv"
)

//work模式和简单模式一样，一条消息只能被一个消费者消费。其和简单模式的区别就在于消费者是否有多个。常用在生产者生产速度大于消费者消费速度的情况下。
func main() {
	//新建一个rabbitMq的实例
	rabbitMq1 := rabbitmq.NewRabbitMQRoute("immocRoute", "immoc.route.one")
	rabbitMq2 := rabbitmq.NewRabbitMQRoute("immocRoute", "immoc.route.two")

	for i := 1; i <= 10; i++ {
		rabbitMq1.RoutePublish("Hello route one! " + strconv.Itoa(i))
		rabbitMq2.RoutePublish("Hello route two! " + strconv.Itoa(i*10))
	}
	fmt.Println("发送成功")
}
