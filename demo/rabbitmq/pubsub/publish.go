package main

import (
	rabbitmq2 "demo/rabbitmq"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmq := rabbitmq2.NewRabbitMQPubSub("newProduct")

	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("PubSub Msg-" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
