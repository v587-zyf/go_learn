package handle

import (
	"comm/t_proto/out/server"
	"context"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"time"
)

const (
	CHAN_SIZE = 1024 * 1024 * 5

	NIL_SLEEP_TIME    = 5 * time.Second
	NO_MSG_SLEEP_TIME = 50 * time.Millisecond
)

type GrpcServer struct {
	gateMsgCh   chan *server.MessageData
	gameMsgCh   chan *server.MessageData
	loginMsgCh  chan *server.MessageData
	centerMsgCh chan *server.MessageData
}

var r = new(GrpcServer)

func init() {
	r.gateMsgCh = make(chan *server.MessageData, CHAN_SIZE)
	r.gameMsgCh = make(chan *server.MessageData, CHAN_SIZE)
	r.loginMsgCh = make(chan *server.MessageData, CHAN_SIZE)
	r.centerMsgCh = make(chan *server.MessageData, CHAN_SIZE)
}

func GetGrpcServer() *GrpcServer {
	return r
}

func (r *GrpcServer) Register(context.Context, *server.RegisterReq) (*server.RegisterAck, error) {

	return nil, nil
}

func (r *GrpcServer) TurnMsg(msg *server.MessageData) {
	switch msg.GetReceiver() {
	case enums.SERVER_GATE:
		r.SendMsgToGate(msg)
	case enums.SERVER_GAME:
		r.SendMsgToGame(msg)
	case enums.SERVER_LOGIN:
		r.SendMsgToLogin(msg)
	case enums.SERVER_CENTER:
		r.SendMsgToCenter(msg)
	default:
		log.Error("unknown receive", zap.Int32("receive", msg.GetReceiver()))
	}
}
