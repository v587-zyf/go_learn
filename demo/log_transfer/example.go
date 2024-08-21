package log_transfer

import (
	"demo/log_transfer/conf"
	"demo/log_transfer/es"
	"demo/log_transfer/kafka"
	"fmt"
	"gopkg.in/ini.v1"
)

/**
 * 将日志数据从kafka取出发到ES
 */

var (
	cfg = new(conf.LogTransferConf)
)

func LogTransferDo() {
	if err := ini.MapTo(cfg, "./log_transfer/conf/conf.ini"); err != nil {
		fmt.Println("config load err:", err)
		return
	}
	fmt.Println("config load success ", cfg)

	if err := es.Init(cfg.EsConf.Address); err != nil {
		fmt.Println("es init err:", err)
		return
	}
	fmt.Println("es init success")

	addrs := []string{
		cfg.KafkaConf.Address,
	}
	if err := kafka.Init(addrs, cfg.KafkaConf.Topic); err != nil {
		fmt.Println("kafka init err:", err)
		return
	}
	fmt.Println("kafka init success")

	select {}
}
