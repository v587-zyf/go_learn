package global

import (
	"go.uber.org/zap"
	"kernel/log"
	"os/signal"
	"syscall"
)

func Run() {
	// server
	{
		log.Info("[HttpServer Start]", zap.String("addr", Conf.Http.ListenAddr))
		go HttpServer.Start()
	}

	go func() {
		<-signalChan

		log.Info("[Server Stop]")

		HttpServer.Stop()

		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
