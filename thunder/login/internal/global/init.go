package global

import (
	"comm/t_tdb"
	"context"
	"github.com/v587-zyf/gc/db/mongo"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
	"github.com/v587-zyf/gc/gcnet/http_server"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/v587-zyf/gc/telegram/go_tg_bot"
	"github.com/v587-zyf/gc/utils"
	"github.com/v587-zyf/gc/workerpool"
	"go.uber.org/zap"
	"login/internal/module/grpc"
	"login/internal/module/handle"
	"login/internal/module/telegram"
	"math/rand"
	"time"
)

func Init() (err error) {
	rand.NewSource(time.Now().UnixNano())
	ctx := context.Background()

	// init config
	{
		if err = log.Init(ctx, log.WithSerName("login"), log.WithSkipCaller(2)); err != nil {
			panic("Log Init err" + err.Error())
		}

		Conf = new(Config)
		if err = utils.Load(Conf); err != nil {
			log.Error("load config file failed", zap.String("err", err.Error()))
			return nil
		}

		if err = t_tdb.Init(Conf.Table.Path); err != nil {
			log.Error("table init err", zap.Error(err))
			return
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

		HttpServer = http_server.NewHttpServer()
		if err = HttpServer.Init(ctx, http_server.WithListenAddr(Conf.Http.ListenAddr),
			http_server.WithIsHttps(Conf.Http.Https), http_server.WithKey(Conf.Http.Key),
			http_server.WithPem(Conf.Http.Pem), http_server.WithAllOrigins(Conf.Http.AllowOrigins)); err != nil {
			log.Error("start server failed", zap.String("listenAddr", Conf.Http.ListenAddr), zap.Error(err))
			return
		}
		if err = grpc_client.InitGrpcClient(ctx, grpc_client.WithListenAddr(Conf.Grpc.ListenAddr)); err != nil {
			log.Error("start server failed",
				zap.String("listenAddr", Conf.Grpc.ListenAddr), zap.String("err", err.Error()))
			return
		}
		if err = grpc.GrpcInit(uint64(Snowflake.Generate())); err != nil {
			log.Error("grpc init err", zap.Error(err))
			return err
		}
		if err = handle.HandleInit(handle.WithTgLoginToken(Conf.Telegram.LoginToken), handle.WithTgHttpServer(HttpServer),
			handle.WithTgClientUrl(Conf.Telegram.ClientUrl), handle.WithTgStartPhoto(Conf.Telegram.StartPhoto),
			handle.WithTgStartCaption(Conf.Telegram.StartCaption), handle.WithSID(Conf.SID)); err != nil {
			log.Error("handle init err", zap.Error(err))
			return err
		}
	}

	// init tg bot
	{
		if !Conf.Dev {
			if err = go_tg_bot.InitTgBot(ctx, go_tg_bot.WithToken(Conf.Telegram.LoginToken)); err != nil {
				log.Error("tg bot init err", zap.Error(err))
				return err
			}
			telegram.TelegramInit()
		}
	}

	return nil
}
