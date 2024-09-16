package main

import (
	"core/log"
	"go.uber.org/zap"
	"login/internal/global"
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
