package handle

import (
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
	"github.com/v587-zyf/gc/gcnet/grpc_msg"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/workerpool"
	"go.uber.org/zap"
	"io"
	"kernel/handler"
	"time"
)

var (
	serSfId uint64
	stream  server.RegisterService_GameMsgClient
)

func GrpcInit(id uint64) {
	var err error

	serSfId = id
	stream, err = server.NewRegisterServiceClient(grpc_client.GetClient()).GameMsg(grpc_client.GetCtx())
	if err != nil {
		panic(err)
		return
	}

	regReq := &server.RegisterReq{Id: uint64(serSfId)}
	reqBytes, err := handler.GetClientWsHandler().Marshal(server.MsgID_Register_ReqId, 0, 0, regReq)
	if err != nil {
		panic(err)
		return
	}

	msgData := &server.MessageData{Sender: 0, Receiver: 0, Content: reqBytes}
	if err = stream.Send(msgData); err != nil {
		log.Error("grpc send err", zap.Error(err))
		return
	}
}

func GrpcListen() {
	var err error

	go func() {
	LOOP:
		for {
			select {
			case msg := <-grpc_msg.GetMsg():
				//log.Debug("1------------------send to game", zap.Any("msg", msg))
				if err = stream.Send(msg.(*server.MessageData)); err != nil {
					log.Error("send to game err", zap.Error(err))
					break LOOP
				}
			default:
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()

	var (
		msg      *server.MessageData
		frameMsg *iface.MessageFrame
	)

	for {
		msg, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("grpc route recv err", zap.Error(err))
			break
		}

		switch msg.GetMsgType() {
		case server.MsgType_Client:
			frameMsg, err = handler.GetClientWsHandler().UnmarshalClient(msg.GetContent())
		case server.MsgType_Server:
			frameMsg, err = handler.GetClientWsHandler().UnmarshalServer(msg.GetContent())
		default:
			log.Error("grpc route msg type err", zap.Any("msgType", msg.GetMsgType()))
			continue
		}
		if err != nil {
			log.Error("grpc route unmarshal err", zap.Error(err))
			continue
		}

		//log.Debug("2----------------", zap.Any("msg", msg), zap.Any("frameMsg", frameMsg))
		workerpool.AssignWsTask(handler.GetClientWsHandler().GetHandler(frameMsg.MsgID), nil, frameMsg)
	}
}
