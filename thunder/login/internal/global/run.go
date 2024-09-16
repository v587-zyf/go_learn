package global

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"login/internal/module/grpc"
	"os/signal"
	"syscall"
)

func Run() {
	// server
	{
		log.Info("[Grpc Start]", zap.String("addr", Conf.Grpc.ListenAddr))
		go grpc.GrpcListen()

		log.Info("[HttpServer Start]", zap.String("addr", Conf.Http.ListenAddr))
		go HttpServer.Start()

		go grpc.ListenSendGrpc()
	}

	go func() {
		<-signalChan

		log.Info("[Server Stop]")

		HttpServer.Stop()
		grpc.GetGrpc().Stop()

		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
