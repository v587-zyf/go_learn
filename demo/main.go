package main

import (
	"demo/elasticSearch"
	"demo/etcd"
	"demo/grpc"
	"demo/kafka"
	"demo/log_transfer"
	"demo/logagent"
	"demo/telegram"
	"flag"
	"fmt"
)

var (
	mod string
)

func main() {
	flag.StringVar(&mod, "m", "tg", "mod")
	flag.Parse()
	//// 启动所有cpu
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//
	//// 随机种子(精确到纳秒)
	//rand.New(rand.NewSource(time.Now().UnixNano()))
	//
	//M := manager.Get()
	//if err := M.Init(""); err != nil {
	//	fmt.Println("init err:", err)
	//	return
	//}
	//
	//if err := M.Start(); err != nil {
	//	fmt.Println("start err:", err)
	//	return
	//}
	//
	// tool.Do()
	//rabbitmq.Do()

	//udp.Server()
	//udp.Client()

	//tool.DoChan()
	//menu.Do()
	//ticket.Do()
	////example.Time()
	//
	//tools.WaitForTerminate()
	//M.Stop()

	//mysql.DriverDo()
	//mysql.SqlxDo()

	switch mod {
	case "p":
		//nsq.ProducerDo()
		kafka.ProducerDo()
	case "c":
		//nsq.ConsumerDo()
		kafka.ConsumerDo()
	case "e":
		etcd.EtcdDo()
	case "l":
		logagent.Do()
	case "es":
		elasticSearch.Do()
	case "ls":
		log_transfer.LogTransferDo()
	case "tg":
		telegram.TelegramDo()
	case "gs":
		grpc.ServiceDo()
	case "gc":
		grpc.ClientDo()
	default:
		fmt.Println("mod error")
	}
}
