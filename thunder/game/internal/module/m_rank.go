package module

import (
	"comm/t_data/redis"
	pb "comm/t_proto/out/client"
	"errors"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
	"strconv"
)

type RankMgr struct {
	module.DefModule
}

func NewRankMgr() *RankMgr {
	return &RankMgr{}
}

func (s *RankMgr) Rank(userID uint64, req *pb.RankReq) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}
	locker.Unlock()

	var me *pb.RankUnit
	if rdbUser.Basic.Lv == int(req.GetLv()) {
		meRanking, err := redis.GetUserRankingByType(req.GetRankType(), userID, int(req.GetLv()))
		if err != nil {
			log.Error("get me ranking err", zap.Error(err))
			return nil, 0, err
		}
		meScore, err := redis.GetUserRankScoreByType(req.GetRankType(), userID, int(req.GetLv()))
		if err != nil {
			log.Error("get me score err", zap.Error(err))
			return nil, 0, err
		}
		me = &pb.RankUnit{
			Uid:       userID,
			Ranking:   int32(meRanking + 1),
			Head:      rdbUser.Telegram.Head,
			FirstName: rdbUser.Telegram.FirstName,
			LastName:  rdbUser.Telegram.LastName,
			UserName:  rdbUser.Telegram.UserName,
			Gold:      meScore,
		}
	}

	var rankInfoSlice []*pb.RankUnit
	rdbRankDatas, err := redis.GetRankByType(req.GetRankType(), int(req.GetLv()))
	if err != nil && !errors.Is(err, redis.NIL) {
		err = errcode.ERR_PARAM
		log.Error("redis.GetRankListByType err", zap.Error(err))
		return
	}

	var (
		rankRdbUser   *redis.User
		rankRdbLocker *rdb_cluster.Locker

		unlockFN = func() {}
	)
	rankInfoSlice = make([]*pb.RankUnit, len(rdbRankDatas))
	for k, v := range rdbRankDatas {
		rankUID, err := strconv.ParseUint(v.Member.(string), 10, 64)
		if err != nil {
			log.Error("strconv.ParseUint err", zap.Error(err))
			continue
		}
		rankRdbLocker = redis.LockUser(rankUID, GetClientModuleMgrOptions().SID)
		if rankRdbLocker == nil {
			log.Error("get locker err", zap.Uint64("rankUID", rankUID))
			continue
		}
		unlockFN = func() { rankRdbLocker.Unlock() }
		rankRdbUser, err = redis.GetUser(rankUID, rankRdbLocker)
		if err != nil {
			unlockFN()
			log.Error("get redis user err", zap.Error(err), zap.Uint64("rankUID", rankUID))
			continue
		}

		rankInfoSlice[k] = &pb.RankUnit{
			Uid:       rankUID,
			Ranking:   int32(k) + 1,
			Head:      rankRdbUser.Telegram.Head,
			FirstName: rankRdbUser.Telegram.FirstName,
			LastName:  rankRdbUser.Telegram.LastName,
			UserName:  rankRdbUser.Telegram.UserName,
			Gold:      rankRdbUser.Basic.Gold,
		}
		unlockFN()
	}

	msgID = pb.MsgID_Rank_AckId
	msg = &pb.RankAck{
		RankType: req.GetRankType(),
		Lv:       req.GetLv(),
		List:     rankInfoSlice,
		Me:       me,
	}

	return
}
