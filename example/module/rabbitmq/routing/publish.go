package main

import (
	"example/module/rabbitmq"
	"fmt"
	"strconv"
	"time"
)

func main() {
	rabbitmqOne := rabbitmq.NewRabbitMQRouting("Routing", "one")
	rabbitmqTwo := rabbitmq.NewRabbitMQRouting("Routing", "two")

	for i := 0; i < 100; i++ {
		rabbitmqOne.PublishRouting("Routing one Msg-" + strconv.Itoa(i))
		rabbitmqTwo.PublishRouting("Routing two Msg-" + strconv.Itoa(i))
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
