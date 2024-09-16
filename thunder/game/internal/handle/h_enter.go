package handle

import (
	"comm/comm"
	"comm/t_enum"
	"comm/t_errcode"
	pb "comm/t_proto/out/client"
	"comm/t_proto/out/server"
	"game/internal/module"
	"github.com/v587-zyf/gc/iface"
	"kernel/handler"
)

func init() {
	handler.GetClientWsHandler().Register(server.MsgID_UserOnline_NtfId, Enter)
	handler.GetClientWsHandler().Register(pb.MsgID_Gm_ReqId, Gm)
}

func Enter(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_ENTER).(*module.EnterMgr).Enter(userID)
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func Gm(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	var err error
	if !GetHandleOps().Dev {
		err = errCode.ERR_GM_CLOSE
		comm.SendErr2User(userID, err)
		return
	}

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_ENTER).(*module.EnterMgr).GM(userID, msgFrame.Body.(*pb.GmReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
