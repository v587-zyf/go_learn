package module

import (
	"comm/t_data/redis"
	"comm/t_errcode"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"go.uber.org/zap"
)

type ShopMgr struct {
	module.DefModule
}

func NewShopMgr() *ShopMgr {
	return &ShopMgr{}
}

func (s *ShopMgr) ShopBuy(userID uint64, req *pb.ShopBuyReq) (msg iface.IProtoMessage, msgID int32, err error) {
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

	shopCfg := t_tdb.GetShopShopCfg(int(req.GetShopId()))
	if shopCfg == nil {
		err = errCode.ERR_CONF_NIL
		return
	}
	shopType := pb.ShopType(shopCfg.Type)

	nowNum := rdbUser.Shop.GetNum(shopType)
	if nowNum >= shopCfg.Buy_num {
		err = errCode.ERR_BUY_NUM_MAX
		return
	}
	if rdbUser.Basic.Diamond < shopCfg.Diamond {
		err = errCode.ERR_DIAMOND_NOT_ENOUGH
		return
	}

	rdbUser.LessDiamond(shopCfg.Diamond)
	rdbUser.Shop.AddNum(shopType)

	switch shopType {
	case pb.ShopType_Normal_Strength:
		userCfg := t_tdb.GetUserUserCfg(rdbUser.Basic.Lv)
		if userCfg == nil {
			log.Error("user conf not found", zap.Uint64("userID", userID), zap.Int("userLv", rdbUser.Basic.Lv))
			err = errCode.ERR_CONF_NIL
			return
		}
		rdbUser.Basic.Strength = userCfg.Strength_max
		rdbUser.Basic.ExtraStrength = rdbUser.Shop.GetAllShopStrength()
	case pb.ShopType_Extra_Strength:
		//rdbUser.Basic.ExtraStrength += t_tdb.Conf().Buy_extra_strength
	}

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("save user fail", zap.Uint64("userID", userID), zap.Error(err))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_ShopBuy_AckId
	msg = &pb.ShopBuyAck{
		ShopId:        req.GetShopId(),
		Diamond:       int32(rdbUser.Basic.Diamond),
		Strength:      int32(rdbUser.Basic.Strength),
		ExtraStrength: int32(rdbUser.Basic.ExtraStrength),
		Info:          rdbUser.Shop.ToPbByType(shopType),
	}

	return
}
