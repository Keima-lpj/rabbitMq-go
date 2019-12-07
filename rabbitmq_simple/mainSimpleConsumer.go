package rabbitmq_simple

import "immoc-rabbitmq/rabbitmq"

func main() {
	rabbitMq := rabbitmq.NewRabbitMQSimple("immocSimple")
	rabbitMq.SimpleConsume()
}
