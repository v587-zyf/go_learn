package global

import (
	"comm/t_data/redis"
	"gate/internal/handle"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"os/signal"
	"syscall"
	"time"
)

var (
	ch = make(chan struct{}, 1)
)

func Run() {
	// server
	{
		log.Info("[Grpc Start]", zap.String("addr", Conf.Grpc.ListenAddr))
		go handle.GrpcListen()

		log.Info("[WsServer Start]", zap.String("addr", Conf.Ws.Addr))
		go WsServer.Start()

		go Sync()
	}

	go func() {
		<-signalChan

		ch <- struct{}{}
		log.Info("[Server Stop]")

		grpc_client.Stop()
		WsServer.Stop()
		ws_session.GetSessionMgr().Close()

		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}

func Sync() {
	var setGate = func() {
		redis.SetGate(float64(ws_session.GetSessionMgr().Length()), redis.FormatGateData(Conf.Ws.LinkAddr, int32(Conf.SID)))
	}

	setGate()
	syncTicker := time.NewTicker(time.Second * 5)
	defer syncTicker.Stop()

LOOP:
	for {
		select {
		case <-syncTicker.C:
			setGate()
		case <-ch:
			redis.DelGate(redis.FormatGateData(Conf.Ws.LinkAddr, int32(Conf.SID)))
			break LOOP
		}
	}
}
