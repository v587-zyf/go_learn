package handle

import (
	"comm/comm"
	enum "comm/t_enum"
	pb "comm/t_proto/out/client"
	"game/internal/module"
	"github.com/v587-zyf/gc/iface"
	"kernel/handler"
)

func init() {
	handler.GetClientWsHandler().Register(pb.MsgID_GuildList_ReqId, GuildList)
	handler.GetClientWsHandler().Register(pb.MsgID_GuildRank_ReqId, GuildRank)
	handler.GetClientWsHandler().Register(pb.MsgID_GuildJoin_ReqId, GuildJoin)
	handler.GetClientWsHandler().Register(pb.MsgID_GuildLeave_ReqId, GuildLeave)
}

func GuildList(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_GUILD).(*module.GuildMgr).GuildList(msgFrame.Body.(*pb.GuildListReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func GuildRank(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_GUILD).(*module.GuildMgr).GuildRank(msgFrame.Body.(*pb.GuildRankReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func GuildJoin(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_GUILD).(*module.GuildMgr).GuildJoin(userID, msgFrame.Body.(*pb.GuildJoinReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func GuildLeave(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_GUILD).(*module.GuildMgr).GuildLeave(userID)
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
