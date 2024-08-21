package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

/**
 * sarama 1.19版本后添加了ztcd压缩算法，需要用到cgo
 * 在go.mod添加 github.com/Shopify/sarama v1.19.0
 * 执行 go mod download
 */

type logData struct {
	topic string
	data  string
}

var (
	client      sarama.SyncProducer
	logDataChan chan *logData
)

func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	// 发送完需要leader和follower都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付消息在success_channel返回
	config.Producer.Return.Successes = true

	// 连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}

	logDataChan = make(chan *logData, maxSize)

	go SendToKafka()
	return
}

func SendToChan(topic, data string) {
	msg := &logData{topic: topic, data: data}
	logDataChan <- msg
}

func SendToKafka() {
	for {
		select {
		case ld := <-logDataChan:
			// make msg
			msg := new(sarama.ProducerMessage)
			msg.Topic = ld.topic
			msg.Value = sarama.StringEncoder(ld.data)
			// send msg
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("pid:%d, offset:%d msg:%v\n",
				pid, offset, msg.Value)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

//// example
//func KafkaDo() {
//	config := sarama.NewConfig()
//	// 发送完需要leader和follower都确认
//	config.Producer.RequiredAcks = sarama.WaitForAll
//	// 新选出一个partition
//	config.Producer.Partitioner = sarama.NewRandomPartitioner
//	// 成功交付消息在success_channel返回
//	config.Producer.Return.Successes = true
//
//	// 消息
//	msg := &sarama.ProducerMessage{}
//	msg.Topic = "web_log"
//	msg.Value = sarama.StringEncoder("this is test log")
//
//	// 连接kafka
//	addrs := []string{"127.0.0.1:9092"}
//	client, err := sarama.NewSyncProducer(addrs, config)
//	if err != nil {
//		fmt.Println("producer closed, err:", err)
//		return
//	}
//	fmt.Println("kafka connection success")
//	defer client.Close()
//
//	pid, offset, err := client.SendMessage(msg)
//	if err != nil {
//		fmt.Println("send message failed, err:", err)
//		return
//	}
//	fmt.Printf("pid:%d, offset:%d\n", pid, offset)
//	fmt.Println("send message success")
//}
