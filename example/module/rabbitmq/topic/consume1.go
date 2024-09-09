package main

import "example/module/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMQTopic("Topic", "#")
	rabbitmq.ReceiveTopic()
}
