package main

import "immoc-rabbitmq/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQRoute("immocRoute", "immoc.route.two")
	rabbitMq.RouteConsumer()
}
