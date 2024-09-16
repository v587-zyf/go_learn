package robot

import (
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) ShopBuy() (err error) {
	shopId := InputInt32("请输入shopId")

	req := &pb.ShopBuyReq{
		ShopId: shopId,
	}
	if err = r.SendMsg(pb.MsgID_ShopBuy_ReqId, req); err != nil {
		log.Error("send move err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) ShopBuyAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.ShopBuyAck)

	log.Debug("shopBuy", zap.Any("msg", msg))
}
