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
	handler.GetClientWsHandler().Register(pb.MsgID_Card_ReqId, Card)
	handler.GetClientWsHandler().Register(pb.MsgID_CardUpLv_ReqId, CardLvUp)
}

func Card(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_CARD).(*module.CardMgr).Card(userID)
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
func CardLvUp(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_CARD).(*module.CardMgr).UpLv(userID, msgFrame.Body.(*pb.CardUpLvReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
