package main

import (
	rabbitmq2 "demo/rabbitmq"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := rabbitmq2.NewRabbitMQTopic("Topic", "imooc.topic.one")
	rabbitmqTwo := rabbitmq2.NewRabbitMQTopic("Topic", "imooc.topic.two")

	for i := 0; i < 100; i++ {
		rabbitmqOne.PublishTopic("Topic one Msg-" + strconv.Itoa(i))
		rabbitmqTwo.PublishTopic("Topic two Msg-" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
