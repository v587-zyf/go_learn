package main

import rabbitmq2 "demo/rabbitmq"

func main() {
	rabbitmq := rabbitmq2.NewRabbitMQTopic("Topic", "imooc.*.two")
	rabbitmq.ReceiveTopic()
}
