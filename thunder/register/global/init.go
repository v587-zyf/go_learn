package global

import (
	"comm/t_proto/out/server"
	"context"
	"github.com/v587-zyf/gc/gcnet/grpc_server"
	"github.com/v587-zyf/gc/gcnet/grpc_server_stream_mgr"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"github.com/v587-zyf/gc/workerpool"
	"go.uber.org/zap"
	"math/rand"
	"register/handle"
	"time"
)

func Init() (err error) {
	rand.NewSource(time.Now().UnixNano())
	ctx := context.Background()

	// init config
	{
		if err = log.Init(ctx, log.WithSerName("register"), log.WithSkipCaller(2)); err != nil {
			panic("Log Init err" + err.Error())
		}

		Conf = new(Config)
		if err = utils.Load(Conf); err != nil {
			log.Error("load config file failed", zap.String("err", err.Error()))
			return nil
		}
	}

	// init server
	{
		if err = workerpool.Init(ctx); err != nil {
			log.Error("WorkPool new failed")
			return
		}

		GrpcServer = grpc_server.NewGrpcServer()
		if err = GrpcServer.Init(ctx, grpc_server.WithListenAddr(Conf.Grpc.ListenAddr)); err != nil {
			log.Error("start server failed",
				zap.String("listenAddr", Conf.Grpc.ListenAddr), zap.String("err", err.Error()))
			return
		}
		server.RegisterRegisterServiceServer(GrpcServer.GetServer(), handle.GetGrpcServer())
		if err = grpc_server_stream_mgr.InitGrpcClientStream(); err != nil {
			log.Error("init grpc client stream failed", zap.Error(err))
			return
		}

	}

	return nil
}
