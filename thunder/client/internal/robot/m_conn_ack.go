package robot

import (
	"client/internal/enums"
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Heartbeat(conn iface.IWsSession, data any) {
	//msg := data.(*iface.MessageFrame).Body.(*pb.Heartbeat)
	//
	//log.Debug("heartbeat", zap.Any("msg", msg))
}

func (r *Robot) VerifyAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.VerifyAck)

	if msg.ErrNo == errcode.ERR_SUCCEED.Int32() {
		r.SetStatus(enums.STATUS_IN_GAME)
	}

	log.Debug("verify", zap.Any("msg", msg))
}
func (r *Robot) ReconnectAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.ReconnectAck)

	if msg.ErrNo == errcode.ERR_SUCCEED.Int32() {
		r.SetStatus(enums.STATUS_IN_GAME)
	}

	log.Debug("reconnect", zap.Any("msg", msg))
}
func (r *Robot) KickNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.KickNtf)

	r.SetStatus(enums.STATUS_DISCONNECT)
	log.Debug("kick", zap.Any("msg", msg))
}
func (r *Robot) Test(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.TestMsg)

	log.Debug("Test", zap.Any("msg", msg))
}
func (r *Robot) EnterNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.EnterNtf)

	log.Debug("EnterNtf", zap.Any("msg", msg))
}
func (r *Robot) ErrNtf(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.ErrNtf)

	log.Debug("ErrNtf", zap.Int32("errNo", msg.ErrNo), zap.String("errMsg", msg.ErrMsg))
}
func (r *Robot) GmAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.GmAck)

	log.Debug("GmAck", zap.Any("msg", msg))
}
