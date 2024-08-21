package main

import (
	rabbitmq2 "demo/rabbitmq"
	"fmt"
	"strconv"
)

func main() {
	rabbitmq := rabbitmq2.NewRabbitMQSimple("imoocSimple")

	for i := 0; i < 100; i++ {
		rabbitmq.PublishSimple("Hello imooc " + strconv.Itoa(i))
		fmt.Println(i)
	}
}
