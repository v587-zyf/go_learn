package handle

import (
	"comm/comm"
	"comm/t_data/redis"
	"comm/t_enum"
	pb "comm/t_proto/out/client"
	"comm/t_proto/out/server"
	"game/internal/module"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"kernel/handler"
)

func init() {
	handler.GetClientWsHandler().Register(server.MsgID_UserOff_NtfId, UserOff)
	handler.GetClientWsHandler().Register(server.MsgID_UserIncome_NtfId, UserIncome)

	handler.GetClientWsHandler().Register(pb.MsgID_RedPoint_ReqId, RedPoint)
}

func UserOff(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	module.GetClientModuleMgr().GetModule(enum.G_M_USER).(*module.UserMgr).Off(userID)
}

func UserIncome(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	locker := redis.LockUser(userID, GetHandleOps().SID)
	if locker == nil {
		log.Error("get locker err")
		return
	}
	defer locker.Unlock()

	module.GetClientModuleMgr().GetModule(enum.G_M_USER).(*module.UserMgr).UpLv(userID, locker)
}

func RedPoint(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_USER).(*module.UserMgr).RedPoint(userID, msgFrame.Body.(*pb.RedPointReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
