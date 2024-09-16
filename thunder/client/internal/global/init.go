package global

import (
	"client/internal/robot_mgr"
	"context"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/utils"
	"github.com/v587-zyf/gc/workerpool"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func Init() (err error) {
	rand.NewSource(time.Now().UnixNano())
	ctx := context.Background()

	{
		if err = log.Init(ctx, log.WithSerName("client"), log.WithSkipCaller(2)); err != nil {
			panic("Log Init err" + err.Error())
		}

		Conf = new(Config)
		if err = utils.Load(Conf); err != nil {
			log.Error("load config file failed", zap.String("err", err.Error()))
			return nil
		}

		RobotManager = robot_mgr.NewRobotManager()
		if err = RobotManager.Init(robot_mgr.WithAddr(Conf.Http.ListenAddr), robot_mgr.WithDev(Conf.Dev),
			robot_mgr.WithKey(Conf.Http.Key), robot_mgr.WithPem(Conf.Http.Pem)); err != nil {
			log.Error("RobotManager init failed", zap.Error(err))
			return err
		}
		//RobotManager.Start()
	}

	{
		if err = workerpool.Init(ctx); err != nil {
			log.Error("WorkPool new failed")
			return
		}

	}

	return nil
}
