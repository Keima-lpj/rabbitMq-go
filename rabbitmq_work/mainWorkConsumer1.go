package main

import "rabbitMq-go/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQSimple("immocWork")
	rabbitMq.SimpleConsume()
}
