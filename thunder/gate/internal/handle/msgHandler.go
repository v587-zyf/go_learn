package handle

import (
	"comm/t_data/redis"
	pb "comm/t_proto/out/client"
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/gcnet/grpc_msg"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/handler"
	"reflect"
)

type MsgHandler struct{}

var _ MsgHandler

func (h *MsgHandler) Recv(s iface.IWsSession, data any) {
	//log.Debug("2------------------recv msg", zap.Any("data", data))

	ss, ok := s.(*ws_session.Session)
	if !ok {
		log.Error("session type err", zap.String("ssType", reflect.TypeOf(s).String()))
		return
	}

	dataBytes := data.([]byte)
	msg, err := handler.GetClientWsHandler().UnmarshalClient(dataBytes)
	if err != nil {
		log.Error("msg UnmarshalClient", zap.Error(err))
		return
	}
	if msgProtoType := pb.GetMsgProtoType(msg.MsgID); msgProtoType == nil {
		log.Error("msgID not found", zap.Uint16("msgID", msg.MsgID))
		return
	}

	// heartbeat
	if msg.MsgID == pb.MsgID_HeartbeatId {
		//heart := new(pb.Heartbeat)
		//ss.Send(pb.MsgID_HeartbeatId, msg.Tag, msg.UserID, heart)
		return
	}
	//log.Debug("msg recv", zap.Any("msg", msg))
	//log.Debug("---", zap.Uint64("ssID", ss.GetID()))
	switch msg.MsgID {
	case pb.MsgID_Verify_ReqId:
		Verify(ss, msg)
	case pb.MsgID_Reconnect_ReqId:
		Reconnect(ss, msg)
	default:
		msgData := &server.MessageData{
			Sender:   enums.SERVER_GATE,
			Receiver: enums.SERVER_GAME,
			Content:  dataBytes,
		}
		grpc_msg.SendToMsg(msgData)
		//log.Error("no login. send msg failed", zap.Uint16("msgID", msg.MsgID), zap.Any("msg", msg))
	}

	//log.Debug("msg recv", zap.Any("msg", msg))
	//workerpool.AssignNormalTask(GetHandler(msg.MsgID), ss, msg)
}

func (h *MsgHandler) Start(ss iface.IWsSession) {
	//log.Debug("new client connected", zap.String("addr", ss.GetConn().RemoteAddr().String()))
}

func (h *MsgHandler) Stop(ss iface.IWsSession) {
	//log.Debug("client disconnected", zap.Uint64("userID", ss.GetID()),
	//	zap.String("addr", ss.GetConn().RemoteAddr().String()))

	// 未注册的不管
	UserID := ss.GetID()
	if UserID <= 0 {
		return
	}

	// 删除会话
	defer ws_session.GetSessionMgr().Disconnect(UserID)

	// 设置重连标志
	if err := redis.SetReconnect(UserID); err != nil {
		log.Error("set reconnect failed", zap.Error(err))
		return
	}

	msgID, userOffMsg := makeUserOffMsg()
	SendToCenter(UserID, msgID, userOffMsg)
	SendToGame(UserID, msgID, userOffMsg)
}
