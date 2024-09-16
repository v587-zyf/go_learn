package module

import (
	"comm/t_data/db"
	"comm/t_data/redis"
	errCode "comm/t_errcode"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"go.uber.org/zap"
	"time"
)

type HastenMgr struct {
	module.DefModule
}

func NewHastenMgr() *HastenMgr {
	return &HastenMgr{}
}

func (s *HastenMgr) Hasten(userID uint64, req *pb.HastenReq) (msg iface.IProtoMessage, msgID int32, err error) {
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

	var unit *db.HastenUnit
	timeNow := time.Now()
	switch req.GetHastenType() {
	case pb.HastenType_Free:
		unit = rdbUser.Hasten.Get(pb.HastenType_Free)
		if !unit.StartTime.IsZero() && timeNow.Before(unit.EndTime) {
			err = errCode.ERR_HASTEN_CD
			return
		}
		if rdbUser.Hasten.IsMax() {
			err = errCode.ERR_HASTEN_MAX
			return
		}

		unit.StartTime = timeNow
		unit.EndTime = timeNow.Add(time.Second * time.Duration(t_tdb.Conf().Free_hasten_cd_second))
		rdbUser.Hasten.AddEndTime(time.Second * time.Duration(t_tdb.Conf().Free_hasten_second))
	case pb.HastenType_Diamonds:
		if rdbUser.Basic.Diamond < t_tdb.Conf().Pay_hasten_diamond {
			err = errCode.ERR_DIAMOND_NOT_ENOUGH
			return
		}
		if rdbUser.Hasten.IsMax() {
			err = errCode.ERR_HASTEN_MAX
			return
		}

		unit = rdbUser.Hasten.Get(pb.HastenType_Diamonds)
		unit.StartTime = timeNow
		unit.EndTime = timeNow.Add(time.Second * time.Duration(t_tdb.Conf().Pay_hasten_second))

		rdbUser.LessDiamond(t_tdb.Conf().Pay_hasten_diamond)
		rdbUser.Hasten.AddEndTime(time.Second * time.Duration(t_tdb.Conf().Pay_hasten_second))
	case pb.HastenType_Link:
		unit = rdbUser.Hasten.Get(pb.HastenType_Link)

		err = errcode.ERR_PARAM
		return
	default:
		err = errcode.ERR_PARAM
		return
	}

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("update redis user err", zap.Error(err), zap.Uint64("userID", rdbUser.ID))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_Hasten_AckId
	msg = &pb.HastenAck{
		Hasten: &pb.HastenUnit{
			HastenType: req.GetHastenType(),
			StartTime:  unit.StartTime.Unix(),
			EndTime:    unit.EndTime.Unix(),
		},
		EndTime: rdbUser.Hasten.EndTime.Unix(),
	}

	return
}
