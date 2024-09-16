package robot

import (
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) UpLv() (err error) {
	req := &pb.UpLvReq{}
	if err = r.SendMsg(pb.MsgID_UpLv_ReqId, req); err != nil {
		log.Error("send up lv err", zap.String("err", err.Error()))
	}
	return nil
}

func (r *Robot) UpLvAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.UpLvAck)

	log.Debug("upLv", zap.Any("msg", msg))
}

func (r *Robot) StrengthNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.StrengthNtf)

	log.Debug("StrengthNtf", zap.Any("msg", msg))
}

func (r *Robot) IncomeNtf(conn iface.IWsSession, data any) {
	//msg := data.(*iface.MessageFrame).Body.(*pb.IncomeNtf)
	//
	//log.Debug("IncomeNtf", zap.Any("msg", msg))
}

func (r *Robot) DiamondNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.DiamondNtf)

	log.Debug("DiamondNtf", zap.Any("msg", msg))
}

func (r *Robot) RedPointNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.RedPointNtf)

	log.Debug("RedPointNtf", zap.Any("msg", msg))
}
