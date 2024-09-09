package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func ConsumerDo() {
	config := sarama.NewConfig()
	// 发送完需要leader和follower都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付消息在success_channel返回
	config.Producer.Return.Successes = true

	// 连接kafka
	addrs := []string{"127.0.0.1:9092"}
	client, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		fmt.Println("consumer closed, err:", err)
		return
	}
	fmt.Println("consumer init success")

	topic := "1_log"
	// 订阅topic
	partitionList, err := client.Partitions(topic)
	if err != nil {
		fmt.Println("get partitions err:", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList {
		// 订阅partition
		partitionConsumer, err := client.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("consumer partition err:", err)
			return
		}
		defer partitionConsumer.AsyncClose()

		go func(consumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages() {
				fmt.Println("consumer msg:", string(msg.Value))
			}
		}(partitionConsumer)
	}
	select {}
}
