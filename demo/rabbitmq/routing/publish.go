package main

import (
	rabbitmq2 "demo/rabbitmq"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := rabbitmq2.NewRabbitMQRouting("Routing", "one")
	rabbitmqTwo := rabbitmq2.NewRabbitMQRouting("Routing", "two")

	for i := 0; i < 100; i++ {
		rabbitmqOne.PublishRouting("Routing one Msg-" + strconv.Itoa(i))
		rabbitmqTwo.PublishRouting("Routing two Msg-" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
