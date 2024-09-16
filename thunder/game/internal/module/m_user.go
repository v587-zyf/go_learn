package module

import (
	"comm/comm"
	"comm/t_data/redis"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"fmt"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
	"time"
)

type UserMgr struct {
	module.DefModule
}

func NewUserMgr() *UserMgr {
	return &UserMgr{}
}

func (m *UserMgr) UpLv(userID uint64, locker *rdb_cluster.Locker) {
	var (
		msg   iface.IProtoMessage
		msgID int32
		err   error

		userCfg *t_tdb.UserUserCfg
		oldLv   = 0

		inviteAdd = 0
	)

	defer func() {
		if err != nil {
			comm.SendErr2User(userID, err)
			return
		} else if msg != nil {
			comm.Send2User(userID, msgID, msg)
		}
	}()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		return
	}

	oldLv = rdbUser.Basic.Lv

	for {
		if t_tdb.GetUserUserCfg(rdbUser.Basic.Lv+1) == nil {
			break
		}
		userCfg = t_tdb.GetUserUserCfg(rdbUser.Basic.Lv)
		if userCfg == nil {
			err = fmt.Errorf("user cfg nil. lv:%v", rdbUser.Basic.Lv)
			return
		}

		if rdbUser.Basic.Gold < float64(userCfg.Consume) {
			break
		}

		//rdbUser.LessGold(float64(userCfg.Consume))
		rdbUser.Basic.Lv++
		rdbUser.Basic.FreeResetMap = enums.YES

		inviteAdd += t_tdb.GetInviteLvInviteLvCfg(rdbUser.Basic.Lv).Normal
	}

	if oldLv != rdbUser.Basic.Lv {
		if err = redis.SetUser(rdbUser, locker); err != nil {
			log.Error("set user fail", zap.Uint64("userID", userID), zap.Error(err))
			return
		}

		msgID = pb.MsgID_UpLv_AckId
		msg = &pb.UpLvAck{
			Lv:           int32(rdbUser.Basic.Lv),
			Gold:         rdbUser.Basic.Gold,
			FreeResetMap: rdbUser.Basic.FreeResetMap == enums.YES,
		}

		if rdbUser.Invite.Invite != 0 {
			inviteLock := redis.LockUser(rdbUser.Invite.Invite, GetClientModuleMgrOptions().SID)
			if inviteLock != nil {
				inviteRdbUser, ret := GetRdbUser(rdbUser.Invite.Invite, inviteLock)
				if !ret {
					log.Error("invite user is nil", zap.Uint64("inviteUID", rdbUser.Invite.Invite))
				} else {
					inviteRdbUser.AddDiamond(inviteAdd)
					inviteRdbUser.Invite.AddInvitees(rdbUser.ID, inviteAdd)

					if err = redis.SetUser(inviteRdbUser, inviteLock); err != nil {
						log.Error("update invite user err", zap.Error(err), zap.Uint64("inviteUID", rdbUser.Invite.Invite))
					}
				}
				inviteLock.Unlock()
				comm.Send2User(rdbUser.Invite.Invite, pb.MsgID_Diamond_NtfId, &pb.DiamondNtf{Diamond: int32(inviteRdbUser.Basic.Diamond), DiamondType: pb.DiamondType_Invite})
			}
		}

		if err = redis.RankDelOldData(userID, oldLv); err != nil {
			log.Error("rank del old data fail", zap.Uint64("userID", userID), zap.Int("oldLv", oldLv), zap.Error(err))
		}
		if err = redis.AddGold(userID, rdbUser.Basic.Gold, rdbUser.Basic.Lv); err != nil {
			log.Error("redis add gold err", zap.Uint64("userID", userID), zap.Float64("gold", rdbUser.Basic.Gold), zap.Int("lv", rdbUser.Basic.Lv), zap.Error(err))
		}
	}

	return
}

func (m *UserMgr) Off(userID uint64) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		log.Error("lock user fail", zap.Uint64("userID", userID))
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		log.Error("get rdb user fail", zap.Uint64("userID", userID))
		return
	}
	rdbUser.Basic.LastGoldAt = time.Now()
	if err := redis.SetUser(rdbUser, locker); err != nil {
		log.Error("set user fail", zap.Uint64("userID", userID), zap.Error(err))
		return
	}
}

func (m *UserMgr) RedPoint(userID uint64, req *pb.RedPointReq) (msg iface.IProtoMessage, msgID int32, err error) {
	err = errcode.ERR_PARAM

	//locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	//if locker == nil {
	//	err = errcode.ERR_USER_DATA_INVALID
	//	return
	//}
	//defer locker.Unlock()
	//
	//rdbUser, ret := GetRdbUser(userID, locker)
	//if !ret {
	//	err = errcode.ERR_USER_DATA_NOT_FOUND
	//	return
	//}
	//
	//rdbUser.RedPoint.UnLook(req.GetRedPointType())
	//if err = redis.SetUser(rdbUser, locker); err != nil {
	//	log.Error("set user fail", zap.Uint64("userID", userID), zap.Error(err))
	//	err = errcode.ERR_REDIS_UPDATE_USER
	//	return
	//}
	//
	//msgID = pb.MsgID_RedPoint_AckId
	//msg = &pb.RedPointAck{
	//	RedPoint: rdbUser.RedPoint.ToPbByType(pb.RedPointType_Card),
	//}

	return
}
