package main

import "immoc-rabbitmq/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQSimple("immocWork")
	rabbitMq.SimpleConsume()
}
