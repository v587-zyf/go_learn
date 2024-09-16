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

func CenterListen() {
LOOP:
	for {
		select {
		case centerMsg := <-r.centerMsgCh:
			stream := grpc_server_stream_mgr.RandStreamByType(enums.SERVER_CENTER)
			if stream == nil {
				log.Error("center stream nil")
				time.Sleep(NIL_SLEEP_TIME)
				continue
			}

			if err := stream.(server.RegisterService_CenterMsgServer).Send(centerMsg); err != nil {
				log.Error("send to center err", zap.Error(err))
				break LOOP
			}
		default:
			time.Sleep(NO_MSG_SLEEP_TIME)
		}
	}
}

func (r *GrpcServer) SendMsgToCenter(msg *server.MessageData) {
	r.centerMsgCh <- msg
}

func (r *GrpcServer) CenterMsg(stream server.RegisterService_CenterMsgServer) (err error) {
	var id uint64
	var msg *server.MessageData
	var msgFrame *iface.MessageFrame
	for {
		msg, err = stream.Recv()
		if err != nil {
			grpc_server_stream_mgr.Del(enums.SERVER_CENTER, id)
			//log.Error("recv from center err", zap.Error(err))
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
				grpc_server_stream_mgr.Add(enums.SERVER_CENTER, registerReq.GetId(), stream)
				continue
			}
		}

		r.TurnMsg(msg)
	}
}
