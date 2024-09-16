package handle

import (
	"center/internal/module"
	"comm/t_enum"
	"comm/t_proto/out/server"
	"github.com/v587-zyf/gc/iface"
	"kernel/handler"
)

func init() {
	handler.GetClientWsHandler().Register(server.MsgID_UserOnline_NtfId, Enter)
	handler.GetClientWsHandler().Register(server.MsgID_UserOff_NtfId, Kick)
}

func Enter(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	module.GetClientModuleMgr().GetModule(enum.C_M_GOLD).(*module.GoldMgr).Add(userID)
	module.GetClientModuleMgr().GetModule(enum.C_M_STRENGTH).(*module.StrengthMgr).Add(userID)
}

func Kick(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	module.GetClientModuleMgr().GetModule(enum.C_M_GOLD).(*module.GoldMgr).Remove(userID)
	module.GetClientModuleMgr().GetModule(enum.C_M_STRENGTH).(*module.StrengthMgr).Remove(userID)
}
