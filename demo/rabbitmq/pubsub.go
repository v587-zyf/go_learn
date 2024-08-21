package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

//// MQURL amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
//const MQURL = "amqp://imoocuser:imoocuser@localhost:5672/imooc"
//
//type RabbitMQ struct {
//	conn      *amqp.Connection
//	channel   *amqp.Channel
//	QueueName string // 队列名称
//	Exchange  string // 交换机
//	Key       string // key
//	Mqurl     string // 连接信息
//}
//
//// NewRabbitMQ 创建RabbitMQ实例
//func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
//	var err error
//
//	rabbitmq := &RabbitMQ{
//		QueueName: queueName,
//		Exchange:  exchange,
//		Key:       key,
//		Mqurl:     MQURL,
//	}
//
//	// 连接
//	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
//	rabbitmq.failOnErr(err, "创建连接错误!")
//	rabbitmq.channel, err = rabbitmq.conn.Channel()
//	rabbitmq.failOnErr(err, "channel获取失败")
//
//	return rabbitmq
//}
//
//// Destroy 断开connection和channel 不断开会一直占用资源
//func (r *RabbitMQ) Destroy() {
//	r.conn.Close()
//	r.channel.Close()
//}
//
//// 错误处理
//func (r *RabbitMQ) failOnErr(err error, message string) {
//	if err != nil {
//		log.Fatalf("%s:%s", message, err)
//		panic(fmt.Sprintf("%s:%s", message, err))
//	}
//}

// NewRabbitMQPubSub 订阅模式MQ创建
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, "")
}

// PublishPub 生产
func (r *RabbitMQ) PublishPub(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 模式 fanout
		"fanout",
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
		"",
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

func (r *RabbitMQ) ReceiveSub() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		// 模式 fanout
		"fanout",
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
		// pub/sub模式下 key为空
		"",
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
