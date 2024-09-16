package grpc

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"time"
)

func ListenSendGrpc() {
	go func() {
	LOOP:
		for {
			select {
			case msg := <-g.msgCh:
				if err = stream.Send(msg); err != nil {
					log.Error("send to game err", zap.Error(err))
					break LOOP
				}
			case <-g.exitCh:
				break LOOP
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()
}
