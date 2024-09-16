package module

import (
	"comm/comm"
	"comm/t_data/redis"
	"comm/t_errcode"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"fmt"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
	"math/rand"
)

type CardMgr struct {
	module.DefModule
}

func NewCardMgr() *CardMgr {
	return &CardMgr{}
}

func (m *CardMgr) GetCard(mapID int, userID uint64, locker *rdb_cluster.Locker) {
	var (
		msg   iface.IProtoMessage
		msgID int32
		err   error
	)

	defer func() {
		if err != nil {
			log.Error("get card fail", zap.Uint64("userID", userID), zap.Error(err))
			//comm.SendErr2User(userID, err)
			return
		} else if msg != nil {
			comm.Send2User(userID, msgID, msg)
		}
	}()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}
	userCfg := t_tdb.GetUserUserCfg(rdbUser.Basic.Lv)
	if userCfg == nil {
		err = fmt.Errorf("user cfg nil. lv:%v", rdbUser.Basic.Lv)
		return
	}

	getNum := 0
	for i := 0; i < userCfg.Gacha_num; i++ {
		if rdbUser.Card.GetNum() >= 4 || rand.Intn(101) >= 30 {
			getNum++
			rdbUser.Card.ResetNum()
		} else {
			rdbUser.Card.AddNum()
		}
	}

	if getNum > 0 {
		flag := false
		ack := &pb.CardNtf{Data: make(map[int32]int32)}
		for _, id := range t_tdb.RandCard(mapID, getNum) {
			if rdbUser.Card.Has(id) {
				continue
			}
			rdbUser.Card.Add(id)
			ack.Data[int32(id)] = int32(rdbUser.Card.GetLv(id))
			flag = true
		}
		if flag {
			msgID = pb.MsgID_Card_NtfId
			msg = ack
		}
		//rdbUser.RedPoint.UnLook(pb.RedPointType_Card)
		//comm.Send2User(userID, pb.MsgID_RedPoint_NtfId, &pb.RedPointNtf{RedPoint: []*pb.RedPointUnit{rdbUser.RedPoint.ToPbByType(pb.RedPointType_Card)}})
	}

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("set user fail", zap.Uint64("userID", userID), zap.Error(err))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}
}

func (m *CardMgr) Card(userID uint64) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	msgID = pb.MsgID_Card_AckId
	msg = &pb.CardAck{Data: rdbUser.Card.ToPb()}

	return
}

func (m *CardMgr) UpLv(userID uint64, req *pb.CardUpLvReq) (ntf iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	id := int(req.GetId())

	nowLv := rdbUser.Card.GetLv(id)
	if t_tdb.GetCardLv(id, nowLv+1) == nil {
		err = errCode.ERR_LV_MAX
		return
	}
	cardLvCfg := t_tdb.GetCardLv(id, nowLv)
	if cardLvCfg == nil {
		err = errCode.ERR_CONF_NIL
		return
	}

	switch req.GetConsumeType() {
	case pb.ConsumeType_Gold:
		if rdbUser.Basic.Gold < float64(cardLvCfg.Consume_gold) {
			err = errCode.ERR_GOLD_NOT_ENOUGH
			return
		}
		rdbUser.LessGold(float64(cardLvCfg.Consume_gold))
	case pb.ConsumeType_Diamond:
		if rdbUser.Basic.Diamond < cardLvCfg.Consume_diamond {
			err = errCode.ERR_DIAMOND_NOT_ENOUGH
			return
		}
		rdbUser.LessDiamond(cardLvCfg.Consume_diamond)
	default:
		err = errcode.ERR_PARAM
	}

	rdbUser.Card.AddLv(id)
	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("set user fail", zap.Uint64("userID", userID), zap.Error(err))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_CardUpLv_AckId
	ntf = &pb.CardUpLvAck{
		Id: int32(id),
		Lv: int32(nowLv + 1),
	}

	return
}
