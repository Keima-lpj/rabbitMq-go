package main

import (
	"fmt"
	"immoc-rabbitmq/rabbitmq"
	"strconv"
)

//work模式和简单模式一样，一条消息只能被一个消费者消费。其和简单模式的区别就在于消费者是否有多个。常用在生产者生产速度大于消费者消费速度的情况下。
func main() {
	//新建一个rabbitMq的实例
	rabbitMq := rabbitmq.NewRabbitMQSimple("immocWork")
	for i := 1; i <= 100; i++ {
		rabbitMq.SimplePublish("Hello L! " + strconv.Itoa(i))
	}
	fmt.Println("发送成功")
}
