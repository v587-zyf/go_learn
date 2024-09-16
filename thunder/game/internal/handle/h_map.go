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
	handler.GetClientWsHandler().Register(pb.MsgID_Move_ReqId, Move)
	handler.GetClientWsHandler().Register(pb.MsgID_OpenWall_ReqId, OpenWall)
	handler.GetClientWsHandler().Register(pb.MsgID_GetTreasure_ReqId, GetTreasure)
	handler.GetClientWsHandler().Register(pb.MsgID_OpenTreasure_ReqId, OpenTreasure)
	handler.GetClientWsHandler().Register(pb.MsgID_Revive_ReqId, Revive)
	handler.GetClientWsHandler().Register(pb.MsgID_ResetMap_ReqId, ResetMap)
}

func Move(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_MAP).(*module.MapMgr).Move(userID, msgFrame.Body.(*pb.MoveReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func OpenWall(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_MAP).(*module.MapMgr).OpenWall(userID, msgFrame.Body.(*pb.OpenWallReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func GetTreasure(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_MAP).(*module.MapMgr).GetTreasure(userID, msgFrame.Body.(*pb.GetTreasureReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func OpenTreasure(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_MAP).(*module.MapMgr).OpenTreasure(userID, msgFrame.Body.(*pb.OpenTreasureReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func Revive(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_MAP).(*module.MapMgr).Revive(userID, msgFrame.Body.(*pb.ReviveReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}

func ResetMap(conn iface.IWsSession, data any) {
	msgFrame := data.(*iface.MessageFrame)
	userID := msgFrame.UserID

	msg, msgID, err := module.GetClientModuleMgr().GetModule(enum.G_M_MAP).(*module.MapMgr).ResetMap(userID, msgFrame.Body.(*pb.ResetMapReq))
	if err != nil {
		comm.SendErr2User(userID, err)
		return
	}
	comm.Send2User(userID, msgID, msg)
}
