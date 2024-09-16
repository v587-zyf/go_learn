package robot

import (
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/gcnet/ws_session"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (r *Robot) Init() {
	r.InitSend()

	r.InitRecv()
}

func (r *Robot) InitSend() {
	r.RegisterSend(-1, r.Disconnect, "disconnect")
	r.RegisterSend(1, r.Login, "login")
	sendMap := []struct {
		Desc string
		Fn   any
	}{
		{"connect", r.ConnectLs},
		{"verify", r.Verify},
		{"reconnect", r.Reconnect},
		{"register", r.Register},

		{"move", r.Move},
		{"openWall", r.OpenWall},
		{"revive", r.Revive},
		{"resetMap", r.ResetMap},
		{"getTreasure", r.GetTreasure},
		{"openTreasure", r.OpenTreasure},

		{"upLv", r.UpLv},

		{"card", r.Card},
		{"cardUpLv", r.CardUpLv},

		{"shopBuy", r.ShopBuy},

		{"hasten", r.Hasten},

		{"invite", r.Invite},
		{"inviteReward", r.InviteReward},

		{"rank", r.Rank},
	}

	var actionID int32 = 2
	for _, s := range sendMap {
		if err := r.RegisterSend(actionID, s.Fn, s.Desc); err != nil {
			log.Error("register recv err", zap.Error(err))
		}
		actionID++
	}
	r.RegisterSend(998, r.GM, "gm")
	r.RegisterSend(999, r.StopMenu, "restart")
}

func (r *Robot) InitRecv() {
	recvMap := map[uint16]ws_session.Recv{
		pb.MsgID_TestMsgId: r.Test,

		pb.MsgID_HeartbeatId: r.Heartbeat,
		pb.MsgID_Kick_NtfId:  r.KickNtf,
		pb.MsgID_Err_NtfId:   r.ErrNtf,
		pb.MsgID_Gm_AckId:    r.GmAck,

		pb.MsgID_Verify_AckId:    r.VerifyAck,
		pb.MsgID_Reconnect_AckId: r.ReconnectAck,

		pb.MsgID_Enter_NtfId: r.EnterNtf,

		pb.MsgID_Move_AckId:         r.MoveAck,
		pb.MsgID_OpenWall_AckId:     r.OpenWallAck,
		pb.MsgID_Revive_AckId:       r.ReviveAck,
		pb.MsgID_ResetMap_AckId:     r.ResetMapAck,
		pb.MsgID_GetTreasure_AckId:  r.GetTreasureAck,
		pb.MsgID_OpenTreasure_AckId: r.OpenTreasureAck,

		pb.MsgID_UpLv_AckId: r.UpLvAck,

		pb.MsgID_Strength_NtfId: r.StrengthNtf,
		pb.MsgID_Income_NtfId:   r.IncomeNtf,
		pb.MsgID_Diamond_NtfId:  r.DiamondNtf,
		pb.MsgID_RedPoint_NtfId: r.RedPointNtf,

		pb.MsgID_Card_AckId:     r.CardAck,
		pb.MsgID_Card_NtfId:     r.CardNtf,
		pb.MsgID_CardUpLv_AckId: r.CardUpLvAck,

		pb.MsgID_ShopBuy_AckId: r.ShopBuyAck,

		pb.MsgID_Hasten_AckId: r.HastenAck,

		pb.MsgID_Invite_AckId:       r.InviteAck,
		pb.MsgID_InviteReward_AckId: r.InviteRewardAck,
		pb.MsgID_Invite_NtfId:       r.InviteNtf,

		pb.MsgID_Rank_AckId: r.RankAck,
	}

	for msgID, recv := range recvMap {
		if err := r.RegisterRecv(msgID, recv); err != nil {
			log.Error("register recv err", zap.Error(err))
		}
	}
}
