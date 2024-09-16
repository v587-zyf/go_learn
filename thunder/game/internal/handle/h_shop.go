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
	handler.GetClientWsHandler().Register(pb.MsgID_ShopBuy_ReqId, ShopBuy)
}

func ShopBuy(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_SHOP).(*module.ShopMgr).ShopBuy(userID, msgFrame.Body.(*pb.ShopBuyReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
