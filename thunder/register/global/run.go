package global

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"os/signal"
	"register/handle"
	"syscall"
)

func Run() {
	// server
	{
		log.Info("[Grpc Start]", zap.String("addr", Conf.Grpc.ListenAddr))
		go GrpcServer.Start()

		go handle.CenterListen()
		go handle.GameListen()
		go handle.GateListen()
		go handle.LoginListen()
	}

	go func() {
		<-signalChan

		log.Info("[Server Stop]")

		GrpcServer.Stop()

		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
