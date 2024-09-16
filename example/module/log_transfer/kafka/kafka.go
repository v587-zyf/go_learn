package kafka

import (
	"example/module/log_transfer/es"
	"fmt"
	"github.com/Shopify/sarama"
)

var (
	client sarama.Consumer
)

func Init(addrs []string, topic string) (err error) {
	// 连接kafka
	client, err = sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Println("consumer closed, err:", err)
		return
	}
	fmt.Println("consumer init success")

	// 订阅topic
	partitionList, err := client.Partitions(topic)
	if err != nil {
		fmt.Println("get partitions err:", err)
		return
	}
	fmt.Println("topic:", topic, " ", partitionList)
	for partition := range partitionList {
		// 订阅partition
		pc, err := client.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			continue
		}
		//defer pc.AsyncClose()

		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, string(msg.Value))

				ld := &es.LogData{
					Topic: topic,
					Data:  string(msg.Value),
				}
				es.SendToEsChan(ld)
			}
		}(pc)
	}

	return err
}
