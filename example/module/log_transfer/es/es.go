package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"strings"
	"time"
)

var (
	esClient *elastic.Client
	ch       chan *LogData
)

type LogData struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

func Init(address string) (err error) {
	if !strings.HasPrefix(address, "http://") {
		address = "http://" + address
	}
	// 创建Client, 连接ES
	esClient, err = elastic.NewClient(
		// elasticsearch 服务地址，多个服务地址使用逗号分隔
		elastic.SetURL(address),
		// 基于http base auth验证机制的账号和密码
		elastic.SetBasicAuth("elastic", "f2bc5b4031bc45805c3a31ce916b6bea"),
		// 启用gzip压缩
		elastic.SetGzip(true),
		// 设置监控检查时间间隔
		elastic.SetHealthcheckInterval(10*time.Second),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		// 设置嗅探功能
		elastic.SetSniff(false),
	)
	if err != nil {
		fmt.Printf("elastic connection error: %v\n", err)
		return err
	}

	ch = make(chan *LogData, 100000)

	go SendToEs()

	return
}

func SendToEsChan(msg *LogData) {
	ch <- msg
}

func SendToEs() {
	for {
		select {
		case msg := <-ch:
			res, err := esClient.
				Index().
				Index(msg.Topic).
				BodyJson(msg).
				Do(context.Background())
			if err != nil {
				fmt.Println("es send error:", err)
				continue
			}
			fmt.Printf("id:%v index:%v ", res.Id, res.Index)
		}
	}
}
