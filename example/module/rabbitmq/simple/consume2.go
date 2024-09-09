package main

import "example/module/rabbitmq"

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimple("imoocSimple")
	rabbitmq.ConsumeSimple()
}
