package elastic

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"time"
)

var esClient *elastic.Client

func Init(host string) error {
	var err error

	// 创建Client, 连接ES
	esClient, err = elastic.NewClient(
		// elasticsearch 服务地址，多个服务地址使用逗号分隔
		elastic.SetURL(host),
		// 基于http base auth验证机制的账号和密码
		//elastic.SetBasicAuth("user", "secret"),
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
	} else {
		fmt.Println("elastic connection success")
	}

	//info, code, err := esClient.Ping(host).Do(context.Background())
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	//
	//esVersion, err := esClient.ElasticsearchVersion(host)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Elasticsearch version %s\n", esVersion)
	return nil
}

func GetEsClient() *elastic.Client {
	return esClient
}
