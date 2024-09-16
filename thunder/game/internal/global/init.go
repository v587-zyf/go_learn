package global

import (
	"comm/t_tdb"
	"context"
	"game/internal/handle"
	"game/internal/module"
	"github.com/v587-zyf/gc/db/mongo"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
	"github.com/v587-zyf/gc/gcnet/grpc_msg"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/rdb/rdb_single"
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
		if err = log.Init(ctx, log.WithSerName("game"), log.WithSkipCaller(2)); err != nil {
			panic("Log Init err" + err.Error())
		}

		Conf = new(Config)
		if err = utils.Load(Conf); err != nil {
			log.Error("load config file failed", zap.String("err", err.Error()))
			return err
		}
		if err = t_tdb.Init(Conf.Table.Path); err != nil {
			log.Error("load tdb err", zap.Error(err))
			return err
		}

		Snowflake, err = utils.NewSnowflake(Conf.SID)
		if err != nil {
			log.Error("new snowflake error", zap.Error(err))
			return
		}
	}

	// init db
	{
		if err = mongo.Init(ctx, mongo.WithUri(Conf.Mongo.Uri), mongo.WithDb(Conf.Mongo.DB)); err != nil {
			log.Error("mongo init failed", zap.Error(err))
			return err
		}
		if err = rdb_single.InitSingle(ctx, rdb_single.WithAddr(Conf.Redis.Addr), rdb_single.WithPwd(Conf.Redis.Pass)); err != nil {
			log.Error("redis init failed", zap.Error(err))
			return err
		}
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
		if err = grpc_client.InitGrpcClient(ctx, grpc_client.WithListenAddr(Conf.Grpc.ListenAddr)); err != nil {
			log.Error("start server failed",
				zap.String("listenAddr", Conf.Grpc.ListenAddr), zap.String("err", err.Error()))
			return
		}
		if err = handle.HandleInit(handle.WithDev(Conf.Dev), handle.WithSID(Conf.SID)); err != nil {
			log.Error("handle init err", zap.Error(err))
			return err
		}
		handle.GrpcInit(uint64(Snowflake.Generate()))

		if err = grpc_msg.InitGrpcMsg(ctx); err != nil {
			log.Error("init grpc_msg err", zap.Error(err))
		}
	}

	// init module
	{
		if err = module.Init(module.WithSID(Conf.SID)); err != nil {
			return err
		}
		module.Start()
	}

	return nil
}
