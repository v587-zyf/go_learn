package handle

import (
	"comm/comm"
	"comm/t_enum"
	pb "comm/t_proto/out/client"
	"game/internal/module"
	"github.com/v587-zyf/gc/iface"
	"kernel/handler"
)

func init() {
	handler.GetClientWsHandler().Register(pb.MsgID_Invite_ReqId, Invite)
	handler.GetClientWsHandler().Register(pb.MsgID_InviteReward_ReqId, InviteReward)
}

func Invite(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_INVITE).(*module.InviteMgr).Invite(userID)
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
func InviteReward(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_INVITE).(*module.InviteMgr).InviteReward(userID, msgFrame.Body.(*pb.InviteRewardReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
