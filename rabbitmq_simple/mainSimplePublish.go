package main

import (
	"fmt"
	"rabbitMq-go/rabbitmq"
)

func main() {
	//新建一个rabbitMq的实例
	rabbitMq := rabbitmq.NewRabbitMQSimple("immocSimple")
	rabbitMq.SimplePublish("Hello L!")
	fmt.Println("发送成功")
}
