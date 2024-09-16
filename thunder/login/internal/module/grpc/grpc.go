package grpc

import (
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/gcnet/grpc_client"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/handler"
)

var (
	err     error
	serSfId uint64
	stream  server.RegisterService_LoginMsgClient
)

type Grpc struct {
	msgCh  chan *server.MessageData
	exitCh chan struct{}
}

var g = new(Grpc)

func init() {
	g.msgCh = make(chan *server.MessageData, 1024*1024*5)
	g.exitCh = make(chan struct{}, 1)
}

func GetGrpc() *Grpc {
	return g
}

func (g *Grpc) SendToMsg(msg *server.MessageData) {
	g.msgCh <- msg
}
func (g *Grpc) Stop() {
	g.exitCh <- struct{}{}
}

func GrpcInit(id uint64) (err error) {
	serSfId = id
	stream, err = server.NewRegisterServiceClient(grpc_client.GetClient()).LoginMsg(grpc_client.GetCtx())
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

	return
}
