package main

import "immoc-rabbitmq/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQSimple("immocSimple")
	rabbitMq.SimpleConsume()
}
