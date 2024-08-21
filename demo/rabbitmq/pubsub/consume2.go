package main

import rabbitmq2 "demo/rabbitmq"

func main() {
	rabbitmq := rabbitmq2.NewRabbitMQPubSub("newProduct")
	rabbitmq.ReceiveSub()
}
