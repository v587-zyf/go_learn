package main

import rabbitmq2 "demo/rabbitmq"

func main() {
	rabbitmq := rabbitmq2.NewRabbitMQRouting("Routing", "one")
	rabbitmq.ReceiveRouting()
}
