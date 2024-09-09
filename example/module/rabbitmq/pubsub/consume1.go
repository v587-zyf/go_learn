package main

import "example/module/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub("newProduct")
	rabbitmq.ReceiveSub()
}
