package robot

import (
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Invite() (err error) {
	if err = r.SendMsg(pb.MsgID_Invite_ReqId, &pb.InviteReq{}); err != nil {
		log.Error("send invite err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) InviteAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.InviteAck)

	log.Debug("InviteAck", zap.Any("msg", msg))
}

func (r *Robot) InviteReward() (err error) {
	t := InputInt32("please input reward id")
	req := &pb.InviteRewardReq{
		Id: t,
	}
	if err = r.SendMsg(pb.MsgID_InviteReward_ReqId, req); err != nil {
		log.Error("send invite reward err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) InviteRewardAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.InviteRewardAck)

	log.Debug("InviteRewardAck", zap.Any("msg", msg))
}
func (r *Robot) InviteNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.InviteNtf)

	log.Debug("InviteNtf", zap.Any("msg", msg))
}
