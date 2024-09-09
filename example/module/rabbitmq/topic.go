package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

// NewRabbitMQTopic 路由模式 创建实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, routingKey)
}

// PublishTopic 生产
func (r *RabbitMQ) PublishTopic(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 模式 topic
		"topic",
		true,
		false,
		// true表示exchange不能被client推送消息，仅进行exchange和exchange绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2.发送消息
	r.channel.Publish(
		r.Exchange,
		r.Key,
		// 如果为true，根据exchange类型和routekey规则，
		// 如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		// 如果为true，当exchange发送消息到队列后，发现队列上没绑定消费者
		// 会把消息发送给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

// ReceiveTopic 消费
// 要注意 key
// * 用于匹配一个单词 # 用于匹配多个单词(可以是零个)
// 匹配 imooc.* 表示匹配 imooc.hello
// 匹配 imooc.# 表示匹配 imooc.hello.one
func (r *RabbitMQ) ReceiveTopic() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 模式 topic
		"topic",
		true,
		false,
		// true表示exchange不能被client推送消息，仅进行exchange和exchange绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 2.创建队列 【队列名称】不要写
	q, err := r.channel.QueueDeclare(
		// 随机队列名称
		"",
		// 是否具有排他性
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 队列消费是否阻塞
		false,
		// 额外信息
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	// 3.绑定队列到exchange
	err = r.channel.QueueBind(
		q.Name,
		// routing模式使用初始化时的key
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to bind a queue")

	// 2.接收消息
	msgs, err := r.channel.Consume(
		q.Name,
		// 区分多个消费者
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果为true，不能将同一个connection中
		// 发送的消息传递给这个connection的消费者
		false,
		// 队列消费是否阻塞
		false,
		// 额外信息
		nil,
	)
	r.failOnErr(err, "Failed to Consume")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			// 实现逻辑函数
			log.Printf("Received a message: %s", d.Body)
			//fmt.Println(d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
