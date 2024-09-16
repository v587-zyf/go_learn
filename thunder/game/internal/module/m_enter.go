package module

import (
	"comm/comm"
	"comm/t_data/redis"
	"comm/t_enum"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"fmt"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"math"
	"strings"
	"time"
)

type EnterMgr struct {
	module.DefModule
}

func NewEnterMgr() *EnterMgr {
	return &EnterMgr{}
}

func (m *EnterMgr) Enter(userID uint64) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}

	//dbUser, err := db_user.GetUserModel().GetUserById(userID)
	//if err != nil {
	//	err = errcode.ERR_USER_DATA_INVALID
	//	return
	//}
	userInfo := rdbUser

	if userInfo.Map == nil {
		userInfo.Map, err = MadeMap(userInfo.Basic.Lv)
		if err != nil {
			return
		}

		userInfo.Basic.TreasureCnt = userInfo.Map.Treasure
	}

	userInfo.Init()

	// 离线收益
	{
		// 离线收益获得量=总每小时挂机收益值/3600*离线秒数，向上保留5位小数
		offSec := 0.0
		if !userInfo.Basic.LastGoldAt.IsZero() {
			offSec = time.Since(userInfo.Basic.LastGoldAt).Seconds()
		}
		//if offSec >= 36000 {
		//	offSec = 36000
		//}
		allIncome := 0
		for id, lv := range userInfo.Card.Data {
			if cardLvCfg := t_tdb.GetCardLv(id, lv); cardLvCfg != nil {
				allIncome += t_tdb.GetCardLv(id, lv).Income
			} else {
				log.Error("get card lv err", zap.Uint64("userID", userID), zap.Int("cardID", id), zap.Int("lv", lv))
			}
		}
		addIncome := float64(0)
		// 最后收益时间是否大于加速时间
		if !userInfo.Hasten.EndTime.IsZero() && !userInfo.Basic.LastGoldAt.After(userInfo.Hasten.EndTime) {
			// 当前时间是否小于加速时间
			if time.Now().After(userInfo.Hasten.EndTime) {
				// 大于，计算加成秒速
				multiplySec := userInfo.Hasten.EndTime.Sub(userInfo.Basic.LastGoldAt).Seconds()
				multiplyAddIncome := utils.RoundFloat(float64(allIncome)/3600*multiplySec*2, 5)
				addIncome = utils.RoundFloat(float64(allIncome)/3600*(offSec-multiplySec)*2, 5) + multiplyAddIncome
			} else {
				// 小于收益按2倍
				addIncome = utils.RoundFloat(float64(allIncome)/3600*offSec*2, 5)
			}
		} else {
			// 大于，收益正常
			addIncome = utils.RoundFloat(float64(allIncome)/3600*offSec, 5)
		}

		rdbUser.AddGold(addIncome)
		go comm.Send2User(userID, pb.MsgID_Income_NtfId, &pb.IncomeNtf{Gold: addIncome, IncomeType: pb.IncomeType_Off})
	}

	// 离线恢复体力
	{
		maxStrength := t_tdb.GetUserUserCfg(userInfo.Basic.Lv).Strength_max
		if userInfo.Basic.Strength < maxStrength || userInfo.Basic.ExtraStrength < userInfo.Shop.GetAllShopStrength() {
			offSec := time.Since(userInfo.Basic.LastStrengthAt).Seconds()
			addStrength := int(math.Floor(offSec/float64(t_tdb.Conf().Add_strengthen_second))) * t_tdb.Conf().Add_strengthen_num

			strengthNeed := maxStrength - userInfo.Basic.Strength
			if addStrength < strengthNeed {
				userInfo.Basic.Strength += addStrength
			} else {
				userInfo.Basic.Strength += strengthNeed
				if addStrength-strengthNeed > 0 {
					userInfo.Basic.ExtraStrength += addStrength - strengthNeed
					if userInfo.Basic.ExtraStrength > userInfo.Shop.GetAllShopStrength() {
						userInfo.Basic.ExtraStrength = userInfo.Shop.GetAllShopStrength()
					}
				}
			}
		}
	}

	// invite
	{
		if userInfo.Invite.Flag == enums.NO && userInfo.Invite.Invite != 0 {
			userInfo.Invite.Flag = enums.YES
			userInfo.AddDiamond(t_tdb.Conf().Normal_invite_diamond)

			inviteLock := redis.LockUser(userInfo.Invite.Invite, GetClientModuleMgrOptions().SID)
			if inviteLock != nil {
				inviteRdbUser, ret := GetRdbUser(userInfo.Invite.Invite, inviteLock)
				if !ret {
					log.Error("invite user is nil", zap.Uint64("inviteUID", userInfo.Invite.Invite))
				} else {
					inviteRdbUser.AddDiamond(t_tdb.Conf().Normal_invite_diamond)
					inviteRdbUser.Invite.AddInvitees(userInfo.ID, t_tdb.Conf().Normal_invite_diamond)

					if err = redis.SetUser(inviteRdbUser, inviteLock); err != nil {
						log.Error("update invite user err", zap.Error(err), zap.Uint64("inviteUID", userInfo.Invite.Invite))
					}

					go comm.Send2User(userInfo.Invite.Invite, pb.MsgID_Invite_NtfId, &pb.InviteNtf{Invitees: userID})

					go comm.Send2User(userID, pb.MsgID_Diamond_NtfId, &pb.DiamondNtf{Diamond: int32(userInfo.Basic.Diamond), DiamondType: pb.DiamondType_Invite})
					go comm.Send2User(userInfo.Invite.Invite, pb.MsgID_Diamond_NtfId, &pb.DiamondNtf{Diamond: int32(inviteRdbUser.Basic.Diamond), DiamondType: pb.DiamondType_Invite})
				}
				inviteLock.Unlock()
			}
		}
	}

	//if _, err = db_user.GetUserModel().Upsert(dbUser); err != nil {
	//	log.Error("update db user err", zap.Uint64("userID", dbUser.ID), zap.Error(err))
	//	err = errcode.ERR_MONGO_UPSERT
	//	return
	//}
	//if _, err = redis.UpdateUserByDB(dbUser, locker); err != nil {
	//	log.Error("UpdateUserByDB err", zap.Error(err), zap.Uint64("userID", dbUser.ID))
	//	err = errcode.ERR_REDIS_UPDATE_USER
	//	return
	//}
	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("update redis user err", zap.Error(err), zap.Uint64("userID", rdbUser.ID))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}
	go GetClientModuleMgr().GetModule(enum.G_M_USER).(*UserMgr).UpLv(userInfo.ID, locker)

	fmt.Printf(userInfo.Map.String())
	msgID = pb.MsgID_Enter_NtfId
	msg = &pb.EnterNtf{
		Lv:            int32(userInfo.Basic.Lv),
		Gold:          userInfo.Basic.Gold,
		Strength:      int32(userInfo.Basic.Strength),
		Maps:          redis.BuildPbMap(userInfo.Map),
		Dead:          enums.RES_BOOL[userInfo.Basic.Dead],
		FreeResetMap:  userInfo.Basic.FreeResetMap == enums.YES,
		Head:          userInfo.Telegram.Head,
		Diamond:       int32(userInfo.Basic.Diamond),
		Shop:          userInfo.Shop.ToPb(),
		ExtraStrength: int32(userInfo.Basic.ExtraStrength),
		Hasten:        userInfo.Hasten.ToPb(),
	}

	return
}

func (m *EnterMgr) GM(userID uint64, req *pb.GmReq) (msg iface.IProtoMessage, msgID int32, err error) {
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

	ack := &pb.GmAck{GmType: req.GetGmType(), Data: req.GetData()}

	switch req.GetGmType() {
	case pb.GmType_Set_Gold:
		d := utils.StrToFloat(req.GetData())
		if d >= 999999999 {
			d = 999999999
		}
		rdbUser.Basic.Gold = d
		//comm.Send2User(userID, pb.MsgID_Income_NtfId, &pb.IncomeNtf{Gold: d, IncomeType: pb.IncomeType_OnHook})
	case pb.GmType_Set_Strength:
		d := utils.StrToInt(req.GetData())
		if d >= t_tdb.GetUserUserCfg(rdbUser.Basic.Lv).Strength_max {
			d = t_tdb.GetUserUserCfg(rdbUser.Basic.Lv).Strength_max
		}
		rdbUser.Basic.Strength = d
		//comm.Send2User(userID, pb.MsgID_Strength_NtfId, &pb.StrengthNtf{Strength: int32(d)})
	case pb.GmType_Set_User_Lv:
		d := utils.StrToInt(req.GetData())
		maxLv := 1
		for _, l := range t_tdb.GetLvs() {
			if maxLv < l {
				maxLv = l
			}
		}
		if d > maxLv {
			d = maxLv
		} else if d <= 0 {
			d = 1
		}
		rdbUser.Basic.Lv = d
		//comm.Send2User(userID, pb.MsgID_UpLv_AckId, &pb.UpLvAck{Lv: int32(d), Gold: rdbUser.Basic.Gold, FreeResetMap: rdbUser.Basic.FreeResetMap == enums.YES})
	case pb.GmType_Set_Card_Lv:
		d := strings.Split(req.GetData(), "#")
		if len(d) != 2 {
			err = errcode.ERR_PARAM
			return
		}
		id := utils.StrToInt(d[0])
		lv := utils.StrToInt(d[1])
		rdbUser.Card.SetLv(id, lv)
		//comm.Send2User(userID, pb.MsgID_CardUpLv_AckId, &pb.CardUpLvAck{Id: int32(id), Lv: int32(lv)})
	case pb.GmType_Reset_Card:
		rdbUser.Card.Reset()
		//comm.Send2User(userID, pb.MsgID_Card_AckId, &pb.CardAck{Data: rdbUser.Card.ToPb()})
	case pb.GmType_Set_Diamond:
		d := utils.StrToInt(req.GetData())
		if d >= 999999999 {
			d = 999999999
		}
		rdbUser.Basic.Diamond = d
	case pb.GmType_Set_Shop_Num:
		d := strings.Split(req.GetData(), "#")
		if len(d) != 2 {
			err = errcode.ERR_PARAM
			return
		}
		t := utils.StrToInt(d[0])
		num := utils.StrToInt(d[1])
		rdbUser.Shop.SetNum(pb.ShopType(t), num)
		if pb.ShopType(t) == pb.ShopType_Extra_Strength {
			rdbUser.Basic.ExtraStrength = t_tdb.Conf().Buy_extra_strength * num
		}
	case pb.GmType_Reset_Shop:
		rdbUser.Shop.Reset()
		rdbUser.Basic.ExtraStrength = 0
	case pb.GmType_Reset_Hasten:
		rdbUser.Hasten.Reset()
	case pb.GmType_Reset_Invite_Reward:
		rdbUser.Invite.Reset()
	default:
		err = errcode.ERR_PARAM
		return
	}

	if err = redis.SetUser(rdbUser, locker); err != nil {
		err = errcode.ERR_REDIS_UPDATE_USER
		log.Error("set user fail", zap.Uint64("userID", userID), zap.Error(err))
		return
	}

	msgID = pb.MsgID_Gm_AckId
	msg = ack

	return
}
