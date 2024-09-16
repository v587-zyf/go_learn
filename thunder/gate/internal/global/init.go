package global

import (
	"context"
	"gate/internal/handle"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
	"github.com/v587-zyf/gc/gcnet/grpc_msg"
	"github.com/v587-zyf/gc/gcnet/ws_server"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/utils"
	"github.com/v587-zyf/gc/workerpool"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func Init() (err error) {
	rand.NewSource(time.Now().UnixNano())
	ctx := context.Background()

	// init config
	{
		if err = log.Init(ctx, log.WithSerName("gate"), log.WithSkipCaller(2)); err != nil {
			panic("Log Init err" + err.Error())
		}

		Conf = new(Config)
		if err = utils.Load(Conf); err != nil {
			log.Error("load config file failed", zap.String("err", err.Error()))
			return nil
		}
		//log.Debug("config load succ", zap.Any("config", Conf))

		Snowflake, err = utils.NewSnowflake(Conf.SID)
		if err != nil {
			log.Error("new snowflake error", zap.Error(err))
			return
		}
	}

	// init db
	{
		//if err = rdb_single.InitSingle(ctx, rdb_single.WithAddr(Conf.Redis.Addr), rdb_single.WithPwd(Conf.Redis.Pass)); err != nil {
		//	log.Error("redis init failed", zap.Error(err))
		//	return err
		//}
		if err = rdb_cluster.InitCluster(ctx, rdb_cluster.WithAddr(Conf.Redis.Addrs), rdb_cluster.WithPwd(Conf.Redis.Pass)); err != nil {
			log.Error("redis init failed", zap.Error(err))
			return err
		}
	}

	// init server
	{
		if err = workerpool.Init(ctx); err != nil {
			log.Error("WorkPool new failed")
			return
		}

		WsServer = ws_server.NewWsServer()
		if err = WsServer.Init(ctx, ws_server.WithAddr(Conf.Ws.Addr), ws_server.WithPem(Conf.Ws.Pem),
			ws_server.WithKey(Conf.Ws.Key), ws_server.WithMethod(new(handle.MsgHandler)), ws_server.WithHttps(Conf.Ws.Https)); err != nil {
			log.Error("init ws_server failed", zap.Int64("SID", Conf.SID),
				zap.String("listenAddr", Conf.Tcp.ListenAddr), zap.Error(err))
			return
		}

		if err = grpc_client.InitGrpcClient(ctx, grpc_client.WithListenAddr(Conf.Grpc.ListenAddr)); err != nil {
			log.Error("init grpc_client failed",
				zap.String("listenAddr", Conf.Grpc.ListenAddr), zap.Error(err))
			return
		}
		handle.GrpcInit(uint64(Snowflake.Generate()))

		if err = grpc_msg.InitGrpcMsg(ctx); err != nil {
			log.Error("init grpc_msg err", zap.Error(err))
		}
	}

	return nil
}
