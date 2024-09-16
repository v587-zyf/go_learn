package robot

import (
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Move() (err error) {
	x := InputInt32("请输入X")
	y := InputInt32("请输入Y")

	req := &pb.MoveReq{
		X: x,
		Y: y,
	}
	if err = r.SendMsg(pb.MsgID_Move_ReqId, req); err != nil {
		log.Error("send move err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) MoveAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.MoveAck)

	log.Debug("move", zap.Any("msg", msg))
}

func (r *Robot) OpenWall() (err error) {
	x := InputInt32("请输入X")
	y := InputInt32("请输入Y")

	req := &pb.OpenWallReq{
		X: x,
		Y: y,
	}
	if err = r.SendMsg(pb.MsgID_OpenWall_ReqId, req); err != nil {
		log.Error("send open wall err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) OpenWallAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.OpenWallAck)

	log.Debug("open wall", zap.Any("msg", msg))
}

func (r *Robot) Revive() (err error) {
	req := &pb.ReviveReq{}
	if err = r.SendMsg(pb.MsgID_Revive_ReqId, req); err != nil {
		log.Error("send revive err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) ReviveAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.ReviveAck)

	log.Debug("revive", zap.Any("msg", msg))
}

func (r *Robot) ResetMap() (err error) {
	req := &pb.ResetMapReq{}
	if err = r.SendMsg(pb.MsgID_ResetMap_ReqId, req); err != nil {
		log.Error("send resetMap err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) ResetMapAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.ResetMapAck)

	log.Debug("resetMap", zap.Any("msg", msg))
}

func (r *Robot) GetTreasure() (err error) {
	x := InputInt32("请输入X")
	y := InputInt32("请输入Y")
	req := &pb.GetTreasureReq{X: x, Y: y}
	if err = r.SendMsg(pb.MsgID_GetTreasure_ReqId, req); err != nil {
		log.Error("send getTreasure err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) GetTreasureAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.GetTreasureAck)

	log.Debug("GetTreasureAck", zap.Any("msg", msg))
}

func (r *Robot) OpenTreasure() (err error) {
	x := InputInt32("请输入X")
	y := InputInt32("请输入Y")
	req := &pb.OpenTreasureReq{X: x, Y: y}
	if err = r.SendMsg(pb.MsgID_OpenTreasure_ReqId, req); err != nil {
		log.Error("send openTreasure err", zap.String("err", err.Error()))
	}
	return nil
}
func (r *Robot) OpenTreasureAck(conn iface.IWsSession, data any) {
	msg := data.(*iface.MessageFrame).Body.(*pb.OpenTreasureAck)

	log.Debug("OpenTreasureAck", zap.Any("msg", msg))
}
