package main

import "immoc-rabbitmq/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQPubSub("immocPubSub")
	rabbitMq.PubSubConsumer()
}
