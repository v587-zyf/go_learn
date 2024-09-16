package handle

import (
	pb "comm/t_proto/out/client"
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/iface"
	"kernel/handler"
	"reflect"
)

func init() {
	handler.GetClientWsHandler().Register(server.MsgID_Send2UserId, Send2User)
	handler.GetClientWsHandler().Register(server.MsgID_BroadcastId, Send2User)
}

func Send2User(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID
	send2UserMsg := msgFrame.Body.(*server.Send2User)
	if len(send2UserMsg.GetContent()) <= 0 {
		return
	}

	protoType := pb.GetMsgProtoType(uint16(send2UserMsg.MsgID))
	msg := reflect.New(protoType).Interface().(iface.IProtoMessage)
	if err = msg.Unmarshal(send2UserMsg.GetContent()); err != nil {
		return
	}

	ws_session.GetSessionMgr().Once(userID, func(SS iface.IWsSession) {
		if SS != nil {
			SS.Send2User(uint16(send2UserMsg.GetMsgID()), msg)
		}
	})
}

func Broadcast(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	broadcastMsg := msgFrame.Body.(*server.Broadcast)
	if len(broadcastMsg.GetContent()) <= 0 {
		return
	}

	protoType := pb.GetMsgProtoType(uint16(broadcastMsg.MsgID))
	msg := reflect.New(protoType).Interface().(iface.IProtoMessage)
	if err = msg.Unmarshal(broadcastMsg.GetContent()); err != nil {
		return
	}

	ws_session.GetSessionMgr().Range(func(UID uint64, SS iface.IWsSession) {
		if SS != nil {
			SS.Send2User(uint16(broadcastMsg.GetMsgID()), msg)
		}
	})
}
