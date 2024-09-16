package robot

import (
	pb "comm/t_proto/out/client"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Card() (err error) {
	req := &pb.CardReq{}
	if err = r.SendMsg(pb.MsgID_Card_ReqId, req); err != nil {
		log.Error("send card err", zap.Error(err))
	}
	return nil
}
func (r *Robot) CardAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.CardAck)

	log.Debug("CardAck", zap.Any("msg", msg))
}

func (r *Robot) CardUpLv() (err error) {
	cardID := InputInt32("请输入卡牌id")
	consumeTypeStr := "请输入消费类型\n"
	for k, v := range pb.ConsumeType_name {
		consumeTypeStr += fmt.Sprintf("%d->%s\n", k, v)
	}
	consumeType := InputInt32(consumeTypeStr)

	req := &pb.CardUpLvReq{
		Id:          cardID,
		ConsumeType: pb.ConsumeType(consumeType),
	}
	if err = r.SendMsg(pb.MsgID_CardUpLv_ReqId, req); err != nil {
		log.Error("send card up lv err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) CardUpLvAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.CardUpLvAck)

	log.Debug("CardUpLvAck", zap.Any("msg", msg))
}

func (r *Robot) CardNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.CardNtf)

	log.Debug("CardNtf", zap.Any("msg", msg))
}
