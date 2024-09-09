package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func ProducerDo() {
	config := sarama.NewConfig()
	// 发送完需要leader和follower都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付消息在success_channel返回
	config.Producer.Return.Successes = true

	// 消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder("this is test log")

	// 连接kafka
	addrs := []string{"127.0.0.1:9092"}
	client, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	fmt.Println("kafka connection success")
	defer client.Close()

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed, err:", err)
		return
	}
	fmt.Printf("pid:%d, offset:%d\n", pid, offset)
	fmt.Println("send message success")
}
