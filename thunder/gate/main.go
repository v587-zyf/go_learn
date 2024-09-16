package main

import (
	"core/log"
	"gate/internal/global"
	"go.uber.org/zap"
)

func main() {
	var err error

	err = global.Init()
	if err != nil {
		log.Error("init err", zap.Error(err))
		return
	}

	global.Run()
}
