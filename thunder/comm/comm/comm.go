package comm

import (
	pb "comm/t_proto/out/client"
	"comm/t_proto/out/server"
	"errors"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/gcnet/grpc_msg"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"kernel/handler"
)

func Send2User(userID uint64, msgID int32, msg iface.IProtoMessage) {
	msgBytes, err := msg.Marshal()
	if err != nil {
		log.Error("enterNtf marshal err", zap.Error(err))
		return
	}

	content := &server.Send2User{MsgID: msgID, Content: msgBytes}
	reqBytes, err := handler.GetClientWsHandler().Marshal(server.MsgID_Send2UserId, 0, userID, content)
	if err != nil {
		panic(err)
		return
	}

	msgData := &server.MessageData{Sender: enums.SERVER_GAME, Receiver: enums.SERVER_GATE, Content: reqBytes}
	grpc_msg.Get().SendToMsg(msgData)
}

func SendErr2User(userID uint64, err error) {
	errNtf := new(pb.ErrNtf)
	var errCode errcode.ErrCode
	if errors.As(err, &errCode) {
		errNtf.ErrNo = errCode.Int32()
		errNtf.ErrMsg = errCode.Error()
	} else {
		errNtf.ErrNo = errcode.ERR_STANDARD_ERR.Int32()
		errNtf.ErrMsg = err.Error()
	}
	Send2User(userID, pb.MsgID_Err_NtfId, errNtf)
}

func Broadcast(msgID int32, msg iface.IProtoMessage, g *grpc_msg.GrpcMsg) {
	msgBytes, err := msg.Marshal()
	if err != nil {
		log.Error("enterNtf marshal err", zap.Error(err))
		return
	}

	content := &server.Broadcast{MsgID: msgID, Content: msgBytes}
	reqBytes, err := handler.GetClientWsHandler().Marshal(server.MsgID_BroadcastId, 0, 0, content)
	if err != nil {
		panic(err)
		return
	}

	msgData := &server.MessageData{Sender: enums.SERVER_GAME, Receiver: enums.SERVER_GATE, Content: reqBytes}
	g.SendToMsg(msgData)
}
