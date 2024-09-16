package module

import (
	"comm/t_data/db/db_guild"
	"comm/t_data/redis"
	errCode "comm/t_errcode"
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

type GuildMgr struct {
	module.DefModule
}

func NewGuildMgr() *GuildMgr {
	return &GuildMgr{}
}

func (m *GuildMgr) GuildList(req *pb.GuildListReq) (msg iface.IProtoMessage, msgID int32, err error) {
	var rankInfoSlice []*pb.GuildUnit
	rdbRankDatas, err := redis.GetGuildRankByType(req.GetRankType())
	if err != nil && !errors.Is(err, redis.NIL) {
		err = errcode.ERR_PARAM
		log.Error("redis.GetGuildRankByType err", zap.Error(err))
		return
	}

	guildIds := make([]string, len(rdbRankDatas))
	for k, v := range rdbRankDatas {
		guildIds[k] = v.Member.(string)
	}
	rdbGuildDatas, err := redis.GetGuildByIds(guildIds...)
	if err != nil && !errors.Is(err, redis.NIL) {
		err = errcode.ERR_PARAM
		log.Error("redis.GetGuildByIds err", zap.Error(err))
		return
	}

	rankInfoSlice = make([]*pb.GuildUnit, len(rdbRankDatas))
	for k, data := range rdbGuildDatas {
		rankInfoSlice[k] = &pb.GuildUnit{
			GuildId: data.ID,
			Ranking: int32(k) + 1,
			Name:    data.Title,
			Member:  int32(data.Members.Len()),
			Gold:    rdbRankDatas[k].Score,
		}
	}

	msgID = pb.MsgID_GuildList_AckId
	msg = &pb.GuildListAck{
		RankType: req.GetRankType(),
		List:     rankInfoSlice,
	}

	return
}

func (m *GuildMgr) GuildRank(req *pb.GuildRankReq) (msg iface.IProtoMessage, msgID int32, err error) {
	var rankInfoSlice []*pb.RankUnit
	rdbRankDatas, err := redis.GetGuildInRankByType(req.GetRankType(), req.GetGuildId())
	if err != nil && !errors.Is(err, redis.NIL) {
		err = errcode.ERR_PARAM
		log.Error("redis.GetGuildInRankByType err", zap.Error(err))
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

	msgID = pb.MsgID_GuildRank_AckId
	msg = &pb.GuildRankAck{
		RankType: req.GetRankType(),
		GuildId:  req.GetGuildId(),
		List:     rankInfoSlice,
	}

	return
}

func (m *GuildMgr) GuildJoin(userID uint64, req *pb.GuildJoinReq) (msg iface.IProtoMessage, msgID int32, err error) {
	// redis user
	uLock := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if uLock == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer uLock.Unlock()

	rdbUser, ret := GetRdbUser(userID, uLock)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	// check guild join
	if rdbUser.Basic.GuildID != 0 {
		err = errCode.ERR_GUILD_ALREADY_JOIN
		return
	}

	// redis guild
	gLock := redis.LockGuild(rdbUser.Basic.GuildID, GetClientModuleMgrOptions().SID)
	if gLock == nil {
		log.Error("lock guild err")
		err = errCode.ERR_GUILD_DATA_ERR
		return
	}
	defer gLock.Unlock()

	rdbGuild, err := redis.GetGuild(req.GetGuildId(), gLock)
	if err != nil {
		err = errCode.ERR_GUILD_DATA_ERR
		log.Error("get redis guild err", zap.Error(err))
		return
	}

	// mongo guild
	dbGuild, err := db_guild.GetGuildModel().GetGuildById(req.GetGuildId())
	if err != nil {
		err = errCode.ERR_GUILD_DATA_ERR
		log.Error("get mongo guild err", zap.Error(err))
		return
	}

	// mongo guild add user
	dbGuild.Members.Add(userID)
	if _, err = db_guild.GetGuildModel().Upsert(dbGuild); err != nil {
		err = errcode.ERR_MONGO_UPSERT
		log.Error("upsert guild err", zap.Error(err), zap.Any("guild", dbGuild))
		return
	}

	// redis guild add user
	rdbGuild.AddMember(userID)
	if err = redis.SetGuild(rdbGuild, gLock); err != nil {
		err = errCode.ERR_GUILD_UPDATE_ERR
		log.Error("redis update guild err", zap.Error(err))
		return
	}

	// user join guild
	rdbUser.Basic.GuildID = req.GetGuildId()
	if err = redis.SetUser(rdbUser, uLock); err != nil {
		log.Error("save user fail", zap.Uint64("userID", userID), zap.Error(err))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_GuildJoin_AckId
	msg = &pb.GuildJoinAck{
		GuildId: req.GetGuildId(),
	}

	return
}

func (m *GuildMgr) GuildLeave(userID uint64) (msg iface.IProtoMessage, msgID int32, err error) {
	// redis user
	uLock := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if uLock == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer uLock.Unlock()

	rdbUser, ret := GetRdbUser(userID, uLock)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	// check guild join
	if rdbUser.Basic.GuildID == 0 {
		err = errCode.ERR_GUILD_DID_NOT_JOIN
		return
	}

	// redis guild
	gLock := redis.LockGuild(rdbUser.Basic.GuildID, GetClientModuleMgrOptions().SID)
	if gLock == nil {
		log.Error("lock guild err")
		err = errCode.ERR_GUILD_DATA_ERR
		return
	}
	defer gLock.Unlock()

	rdbGuild, err := redis.GetGuild(rdbUser.Basic.GuildID, gLock)
	if err != nil {
		err = errCode.ERR_GUILD_DATA_ERR
		log.Error("get redis guild err", zap.Error(err))
		return
	}

	// mongo guild
	dbGuild, err := db_guild.GetGuildModel().GetGuildById(rdbUser.Basic.GuildID)
	if err != nil {
		err = errCode.ERR_GUILD_DATA_ERR
		log.Error("get mongo guild err", zap.Error(err))
		return
	}

	// mongo guild del user
	dbGuild.Members.Del(userID)
	if _, err = db_guild.GetGuildModel().Upsert(dbGuild); err != nil {
		err = errcode.ERR_MONGO_UPSERT
		log.Error("upsert guild err", zap.Error(err), zap.Any("guild", dbGuild))
		return
	}

	// redis guild del user
	rdbGuild.DelMember(userID)
	if err = redis.SetGuild(rdbGuild, gLock); err != nil {
		err = errCode.ERR_GUILD_UPDATE_ERR
		log.Error("redis update guild err", zap.Error(err))
		return
	}

	// user leave guild
	rdbUser.Basic.GuildID = 0
	if err = redis.SetUser(rdbUser, uLock); err != nil {
		log.Error("save user fail", zap.Uint64("userID", userID), zap.Error(err))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_GuildLeave_AckId
	msg = &pb.GuildLeaveAck{
		Flag: true,
	}

	return
}
