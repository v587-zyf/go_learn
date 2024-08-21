package main

import rabbitmq2 "demo/rabbitmq"

func main() {
	rabbitmq := rabbitmq2.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
}
