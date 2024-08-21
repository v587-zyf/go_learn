package grpc_msg

import (
	"context"
	"kernel/iface"
)

var defGrpcMsg *GrpcMsg

func InitGrpcMsg(ctx context.Context, opts ...any) (err error) {
	defGrpcMsg = NewGrpcMsg()
	if err = defGrpcMsg.Init(ctx, opts...); err != nil {
		return err
	}

	return nil
}

func SendToMsg(msg *server.MessageData) {
	defGrpcMsg.SendToMsg(msg)
}

func GetMsg() <-chan *server.MessageData {
	return defGrpcMsg.GetMsg()
}

func Send2User(userID uint64, msgID int32, msg iface.IProtoMessage) {
	defGrpcMsg.Send2User(userID, msgID, msg)
}

func SendErr2User(userID uint64, err error) {
	defGrpcMsg.SendErr2User(userID, err)
}

func Broadcast(msgID int32, msg iface.IProtoMessage) {
	defGrpcMsg.Broadcast(msgID, msg)
}
