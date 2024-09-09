package main

import (
	"example/module/rabbitmq"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := rabbitmq.NewRabbitMQTopic("Topic", "imooc.topic.one")
	rabbitmqTwo := rabbitmq.NewRabbitMQTopic("Topic", "imooc.topic.two")

	for i := 0; i < 100; i++ {
		rabbitmqOne.PublishTopic("Topic one Msg-" + strconv.Itoa(i))
		rabbitmqTwo.PublishTopic("Topic two Msg-" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
