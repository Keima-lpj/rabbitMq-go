package main

import "rabbitMq-go/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQRoute("immocRoute", "immoc.route.one")
	rabbitMq.RouteConsumer()
}
