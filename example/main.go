package main

import (
	"demo/global"
	"flag"
	"go.uber.org/zap"
	"kernel/log"
)

var (
	mod string
)

func main() {
	flag.StringVar(&mod, "m", "tg", "mod")
	flag.Parse()

	//// 启动所有cpu
	//runtime.GOMAXPROCS(runtime.NumCPU())

	var err error

	err = global.Init()
	if err != nil {
		log.Error("init err", zap.Error(err))
		return
	}

	global.Run()
}
