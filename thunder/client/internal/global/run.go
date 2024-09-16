package global

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
)

func Run() {
	{
		log.Info("[RobotManager Start]", zap.String("addr", Conf.Http.ListenAddr))
		RobotManager.Start()
	}

	go func() {
		<-signalChan

		// todo offline save user data
		log.Info("[Server Stop]")
		RobotManager.Stop()

		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
