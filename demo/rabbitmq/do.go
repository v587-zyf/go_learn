package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

const Mqurls = "amqp://imoocuser:imoocuser@localhost:5672/imooc"

type RabbitMQDemo struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	Mqurl     string
}

func NewRabbitMQDemo(queueName, exchange, key string) *RabbitMQDemo {
	var err error
	r := &RabbitMQDemo{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     Mqurls,
	}

	r.conn, err = amqp.Dial(r.Mqurl)
	if err != nil {
		fmt.Println("amqp dial err:", err)
		return nil
	}
	r.channel, err = r.conn.Channel()
	if err != nil {
		fmt.Println("conn get channel err:", err)
		return nil
	}

	return r
}

func (r *RabbitMQDemo) Destroy() {
	r.conn.Close()
	r.channel.Close()
}

func Do() {

}
