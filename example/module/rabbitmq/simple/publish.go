package main

import (
	"example/module/rabbitmq"
	"fmt"
	"strconv"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQSimple("imoocSimple")

	for i := 0; i < 100; i++ {
		rabbitmq.PublishSimple("Hello imooc " + strconv.Itoa(i))
		fmt.Println(i)
	}
}
