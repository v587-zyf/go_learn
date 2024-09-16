package module

import (
	"comm/comm"
	"comm/t_data/redis"
	errCode "comm/t_errcode"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
)

type InviteMgr struct {
	module.DefModule
}

func NewInviteMgr() *InviteMgr {
	return &InviteMgr{}
}

func (m *InviteMgr) Invite(userID uint64) (msg iface.IProtoMessage, msgID int32, err error) {
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

	var (
		num    = 0
		pbData = make([]*pb.InviteUnit, len(rdbUser.Invite.Invitees))

		inviteRdbUser *redis.User
		inviteLocker  *rdb_cluster.Locker

		unlockFN = func() {}
	)
	for inviteUID, diamond := range rdbUser.Invite.Invitees {
		inviteLocker = redis.LockUser(inviteUID, GetClientModuleMgrOptions().SID)
		if inviteLocker == nil {
			log.Error("get locker err", zap.Uint64("inviteUID", inviteUID))
			continue
		}
		unlockFN = func() { inviteLocker.Unlock() }
		inviteRdbUser, err = redis.GetUser(inviteUID, inviteLocker)
		if err != nil {
			unlockFN()
			log.Error("get redis user err", zap.Error(err), zap.Uint64("inviteUID", inviteUID))
			continue
		}
		pbData[num] = &pb.InviteUnit{
			Uid:       inviteUID,
			Head:      inviteRdbUser.Telegram.Head,
			FirstName: inviteRdbUser.Telegram.FirstName,
			LastName:  inviteRdbUser.Telegram.LastName,
			UserName:  inviteRdbUser.Telegram.UserName,
			Lv:        int32(inviteRdbUser.Basic.Lv),
			Diamond:   int32(diamond),
		}
		num++
		unlockFN()
	}

	msgID = pb.MsgID_Invite_AckId
	msg = &pb.InviteAck{
		Rewards: rdbUser.Invite.ToPbReward(),
		Invites: pbData,
		Flag:    true,
	}

	return
}

func (m *InviteMgr) InviteReward(userID uint64, req *pb.InviteRewardReq) (msg iface.IProtoMessage, msgID int32, err error) {
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
	if _, ok := rdbUser.Invite.Reward[id]; ok {
		err = errCode.ERR_REWARD_ALREADY
		return
	}

	inviteNumCfg := t_tdb.GetInviteNumInviteNumCfg(id)
	if inviteNumCfg == nil {
		log.Error("invite_num cfg not found", zap.Int("id", id))
		err = errCode.ERR_CONF_NIL
		return
	}

	if len(rdbUser.Invite.Invitees) < inviteNumCfg.Num {
		err = errCode.ERR_INVITE_NUM_NOT_ENOUGH
		return
	}

	rdbUser.Invite.AddReward(id)
	rdbUser.Card.Add(inviteNumCfg.Card_id)
	//rdbUser.RedPoint.UnLook(pb.RedPointType_Card)
	//comm.Send2User(userID, pb.MsgID_RedPoint_NtfId, &pb.RedPointNtf{RedPoint: []*pb.RedPointUnit{rdbUser.RedPoint.ToPbByType(pb.RedPointType_Card)}})
	comm.Send2User(userID, pb.MsgID_Card_NtfId, &pb.CardNtf{Data: map[int32]int32{int32(inviteNumCfg.Card_id): 1}})

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("set user fail", zap.Uint64("userID", userID), zap.Error(err))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_InviteReward_AckId
	msg = &pb.InviteRewardAck{
		Id: req.GetId(),
	}
	return
}
