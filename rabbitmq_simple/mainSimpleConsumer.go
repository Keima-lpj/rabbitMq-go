package main

import "rabbitMq-go/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQSimple("immocSimple")
	rabbitMq.SimpleConsume()
}
