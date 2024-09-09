package main

import (
	"example/module/rabbitmq"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := rabbitmq.NewRabbitMQPubSub("newProduct")

	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("PubSub Msg-" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
