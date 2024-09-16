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
	handler.GetClientWsHandler().Register(pb.MsgID_Rank_ReqId, Rank)
}

func Rank(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_RANK).(*module.RankMgr).Rank(userID, msgFrame.Body.(*pb.RankReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
