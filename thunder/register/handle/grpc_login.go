package handle

import (
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/gcnet/grpc_server_stream_mgr"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/handler"
	"time"
)

func LoginListen() {
LOOP:
	for {
		select {
		case loginMsg := <-r.loginMsgCh:
			stream := grpc_server_stream_mgr.RandStreamByType(enums.SERVER_LOGIN)
			if stream == nil {
				log.Error("login stream nil")
				time.Sleep(NIL_SLEEP_TIME)
				continue
			}

			if err := stream.(server.RegisterService_LoginMsgServer).Send(loginMsg); err != nil {
				log.Error("send to login err", zap.Error(err))
				break LOOP
			}
		default:
			time.Sleep(NO_MSG_SLEEP_TIME)
		}
	}
}

func (r *GrpcServer) SendMsgToLogin(msg *server.MessageData) {
	r.loginMsgCh <- msg
}

func (r *GrpcServer) LoginMsg(stream server.RegisterService_LoginMsgServer) (err error) {
	var id uint64
	var msg *server.MessageData
	var msgFrame *iface.MessageFrame
	for {
		msg, err = stream.Recv()
		if err != nil {
			grpc_server_stream_mgr.Del(enums.SERVER_LOGIN, id)
			//log.Error("recv from login err", zap.Error(err))
			return
		}
		if id == 0 {
			msgFrame, err = handler.GetClientWsHandler().UnmarshalServer(msg.GetContent())
			if err != nil {
				log.Error("unmarshal err", zap.Error(err), zap.Any("msg", msg), zap.Any("msgFrame", msgFrame))
				continue
			}
			if msgFrame.MsgID == server.MsgID_Register_ReqId {
				registerReq := msgFrame.Body.(*server.RegisterReq)
				id = registerReq.GetId()
				grpc_server_stream_mgr.Add(enums.SERVER_LOGIN, registerReq.GetId(), stream)
				continue
			}
		}

		r.TurnMsg(msg)
	}
}
