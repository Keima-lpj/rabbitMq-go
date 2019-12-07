package main

import "immoc-rabbitmq/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQRoute("immocRoute", "immoc.route.one")
	rabbitMq.RouteConsumer()
}
