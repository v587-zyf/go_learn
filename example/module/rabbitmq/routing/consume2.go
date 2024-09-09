package main

import "example/module/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMQRouting("Routing", "two")
	rabbitmq.ReceiveRouting()
}
