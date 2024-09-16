package global

import (
	"center/internal/handle"
	"comm/t_data/redis"
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

		go redis.FixUpdateRank()
		go redis.FixUpdateGuildRank()
		go redis.FixUpdateGuildInRank()
	}

	go func() {
		<-signalChan

		// todo offline save user data
		log.Info("[Server Stop]")

		grpc_client.Stop()

		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
