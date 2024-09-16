package main

import (
	"client/internal/global"
	"core/log"
	"go.uber.org/zap"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.New(rand.NewSource(time.Now().UnixNano()))

	var err error

	err = global.Init()
	if err != nil {
		log.Error("init err", zap.Error(err))
		return
	}

	global.Run()
}
