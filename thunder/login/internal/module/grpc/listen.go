package grpc

import (
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/workerpool"
	"go.uber.org/zap"
	"io"
	"kernel/handler"
)

func GrpcListen() {
	var (
		msg      *server.MessageData
		frameMsg *iface.MessageFrame
	)

LOOP:
	for {
		msg, err = stream.Recv()
		if err == io.EOF {
			break LOOP
		}
		if err != nil {
			log.Error("grpc route recv err", zap.Error(err))
			break LOOP
		}

		//log.Debug("1----------------")
		frameMsg, err = handler.GetClientWsHandler().UnmarshalClient(msg.GetContent())

		//log.Debug("2----------------", zap.Any("msg", msg), zap.Any("frameMsg", frameMsg))
		workerpool.AssignWsTask(handler.GetClientWsHandler().GetHandler(frameMsg.MsgID), nil, msg)
	}
}
