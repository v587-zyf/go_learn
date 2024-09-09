package global

import (
	"context"
	"core/telegram/go_tg_bot"
	"demo/module/telegram"
	"go.uber.org/zap"
	"kernel/db/mongo"
	"kernel/log"
	"kernel/server/http_server"
	"kernel/utils"
	"kernel/workerpool"
	"math/rand"
	"time"
)

func Init() (err error) {
	rand.NewSource(time.Now().UnixNano())
	ctx := context.Background()

	// init config
	{
		if err = log.Init(ctx, log.WithSerName("demo"), log.WithSkipCaller(2)); err != nil {
			panic("Log Init err" + err.Error())
		}

		Conf = new(Config)
		if err = utils.Load(Conf); err != nil {
			log.Error("load config file failed", zap.String("err", err.Error()))
			return nil
		}

		Snowflake, err = utils.NewSnowflake(Conf.SID)
		if err != nil {
			log.Error("new snowflake error", zap.Error(err))
			return
		}
	}

	// init db
	{
		if Conf.Mongo.Use {
			if err = mongo.Init(ctx, mongo.WithUri(Conf.Mongo.Uri), mongo.WithDb(Conf.Mongo.DB)); err != nil {
				log.Error("mongo init failed", zap.Error(err))
				return err
			}
		}
		//if err = rdb_single.InitSingle(ctx, rdb_single.WithAddr(Conf.Redis.Addr), rdb_single.WithPwd(Conf.Redis.Pass)); err != nil {
		//	log.Error("redis init failed", zap.Error(err))
		//	return err
		//}
		//if err = rdb_cluster.InitCluster(ctx, rdb_cluster.WithAddr(Conf.Redis.Addrs), rdb_cluster.WithPwd(Conf.Redis.Pass)); err != nil {
		//	log.Error("redis init failed", zap.Error(err))
		//	return err
		//}
	}

	// init server
	{
		if err = workerpool.Init(); err != nil {
			log.Error("WorkPool new failed")
			return
		}

		HttpServer = http_server.NewHttpServer()
		if err = HttpServer.Init(ctx, http_server.WithListenAddr(Conf.Http.ListenAddr),
			http_server.WithIsHttps(false)); err != nil {
			log.Error("start server failed", zap.String("listenAddr", Conf.Http.ListenAddr), zap.Error(err))
			return
		}
	}

	// init tg bot
	{
		if err = go_tg_bot.InitTgBot(ctx, go_tg_bot.WithToken(Conf.Telegram.Token)); err != nil {
			log.Error("tg bot init err", zap.Error(err))
			return err
		}
		telegram.TgInit()
	}

	return nil
}
