package global

import (
	"game/internal/handle"
	"game/internal/module"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

func Run() {
	// server
	{
		log.Info("[Grpc Start]", zap.String("addr", Conf.Grpc.ListenAddr))
		go handle.GrpcListen()

		module.Run()
	}

	go func() {
		<-signalChan

		// todo offline save user data
		log.Info("[Server Stop]")

		grpc_client.Stop()

		module.Stop()

		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
