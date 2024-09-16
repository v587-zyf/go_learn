package logagent

import (
	"example/module/logagent/conf"
	"example/module/logagent/etcd"
	"example/module/logagent/kafka"
	"example/module/logagent/tail_log"
	"fmt"
	"github.com/v587-zyf/gc/utils"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

var (
	cfg = new(conf.AppConf)
)

//func run() {
//	topic := cfg.Topic
//	for {
//		select {
//		// 1.读取日志
//		case line := <-tail_log.ReadChan():
//			// 2.发送到kafka
//			kafka.SendToKafka(topic, line.Text)
//		default:
//			time.Sleep(time.Second)
//		}
//	}
//}

func Do() {
	// 1.init config
	if err := ini.MapTo(cfg, "./logagent/conf/config.ini"); err != nil {
		fmt.Println("config load err:", err)
		return
	}

	// 2.init kafka
	addrs := []string{
		cfg.KafkaConf.Address,
	}
	if err := kafka.Init(addrs, cfg.KafkaConf.ChanMaxSize); err != nil {
		fmt.Println("kafka init err:", err)
		return
	}
	fmt.Println("kafka init success")

	// 3.init etcd
	if err := etcd.Init(cfg.EtcdConf.Address, time.Duration(cfg.EtcdConf.Timeout)*time.Second); err != nil {
		fmt.Println("etcd init err:", err)
		return
	}
	fmt.Println("etcd init success")

	// 3.1.从etcd拉取配置
	ip, err := utils.GetLocalIp()
	if err != nil {
		fmt.Println("get local ip err:", err)
		return
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key, ip)
	//fmt.Println("etcd conf key:", etcdConfKey)
	logEntries, err := etcd.GetConf(etcdConfKey)
	if err != nil {
		fmt.Println("get from etcd conf err:", err)
		return
	}

	tail_log.Init(logEntries)

	// 3.2.哨兵监听配置
	newCConfChan := tail_log.GetNewConfChan()
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, newCConfChan)
	wg.Wait()
}
