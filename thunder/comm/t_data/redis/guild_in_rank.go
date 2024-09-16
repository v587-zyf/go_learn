package redis

import (
	"comm/t_data/db"
	"comm/t_data/db/db_guild_in_rank"
	enum "comm/t_enum"
	pb "comm/t_proto/out/client"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
)

func AddGuildInRankGold(userID, guildID uint64, gold float64) error {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	if err := rc.ZIncrBy(rCtx, FormatGuildInRankDaily(guildID), gold, strconv.Itoa(int(userID))).Err(); err != nil {
		log.Error("add guild in daily rank err", zap.Error(err))
		return err
	}

	if err := rc.ZIncrBy(rCtx, FormatGuildInRankWeekly(guildID), gold, strconv.Itoa(int(userID))).Err(); err != nil {
		log.Error("add guild in weekly rank err", zap.Error(err))
		return err
	}

	return nil
}

func GetGuildInRankByType(rankType pb.RankType, guildID uint64) ([]redis.Z, error) {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	var rankKey string
	switch rankType {
	case pb.RankType_RankDaily:
		rankKey = FormatGuildInRankDaily(guildID)
	case pb.RankType_RankWeekly:
		rankKey = FormatGuildInRankWeekly(guildID)
	default:
		log.Error("get rank by type err", zap.Any("rankType", rankType))
		return nil, errcode.ERR_PARAM
	}

	return rc.ZRevRangeWithScores(rCtx, rankKey, 0, 499).Result()
}

func GuildInRankToDB(rankType pb.RankType, rc *redis.Client, rCtx context.Context, rankKey string, date int, guildId uint64) error {
	rdbRankDatas, err := rc.ZRevRangeWithScores(rCtx, rankKey, 0, 99).Result()
	if err != nil {
		log.Error("Error fetching guild in rank", zap.Error(err))
		return nil
	}
	if len(rdbRankDatas) <= 0 {
		//log.Warn("rank to db data nil", zap.String("rankKey", rankKey))
		return nil
	}

	seqId, err := db.GenGuildInRankIdSeq()
	if err != nil {
		log.Error("gen guild in rank id err", zap.Error(err))
		return err
	}
	addDbData := &db_guild_in_rank.GuildInRank{
		ID:      seqId,
		GuildId: guildId,
		Type:    rankType,
		Date:    date,
	}

	addDbRanks := make([]*db_guild_in_rank.GuildInRankUnit, len(rdbRankDatas))
	for k, v := range rdbRankDatas {
		userId, err := strconv.ParseUint(v.Member.(string), 10, 64)
		if err != nil {
			log.Error("strconv.ParseUint err", zap.Error(err))
			continue
		}
		addDbRanks[k] = &db_guild_in_rank.GuildInRankUnit{
			Ranking: k + 1,
			UserId:  userId,
			Gold:    int(v.Score),
		}
	}
	addDbData.Rank = addDbRanks
	if _, err = db_guild_in_rank.GetGuildInRankModel().Upsert(addDbData); err != nil {
		log.Error("upsert guild in rank err", zap.Error(err))
		return err
	}
	return nil
}

func GuildRankInDel(rc *redis.Client, rCtx context.Context, rankKey ...string) {
	if _, err := rc.Del(rCtx, rankKey...).Result(); err != nil {
		log.Error("del guild in rank err", zap.Error(err))
	}
}
func UpdateGuildInRankByRankType(rankType pb.RankType) error {
	var (
		guildId  uint64
		date     int
		pattern  string
		rankKey  string
		keySlice []string

		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()

		dbRank *db_guild_in_rank.GuildInRank

		delKeys = make([]string, 0)
	)

	switch rankType {
	case pb.RankType_RankDaily:
		pattern = enum.REDIS_KEY_GUILD_IN_RANK_DAILY
	case pb.RankType_RankWeekly:
		pattern = enum.REDIS_KEY_GUILD_IN_RANK_WEEKLY
	default:
		log.Error("get guild in rank by type err", zap.Any("rankType", rankType))
		return errors.New("get guild in rank by type err")
	}
	pattern = pattern + "*"
	keys, err := rc.Keys(rCtx, pattern).Result()
	if err != nil && !errors.Is(err, NIL) {
		log.Error("get guild in rank keys err", zap.Error(err), zap.String("pattern", pattern))
		return err
	}
	if len(keys) <= 0 {
		return nil
	}

	for _, key := range keys {
		keySlice = strings.Split(key, ":")

		guildId, err = strconv.ParseUint(keySlice[1], 10, 64)
		if err != nil {
			log.Error("rank date err", zap.String("key", key), zap.Error(err))
			continue
		}

		date, err = strconv.Atoi(keySlice[2])
		if err != nil {
			log.Error("rank date err", zap.String("key", key), zap.Error(err))
			continue
		}

		switch rankType {
		case pb.RankType_RankDaily:
			rankKey = FormatGuildInRankDaily(guildId, date)
		case pb.RankType_RankWeekly:
			rankKey = FormatGuildInRankWeekly(guildId, date)
		}
		dbRank, err = db_guild_in_rank.GetGuildInRankModel().GetRankInfoByTypeAndDateAndGid(rankType, date, guildId)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			log.Error("get guild in rank info err", zap.Error(err))
			continue
		}
		if dbRank == nil || dbRank.ID == 0 {
			if err = GuildInRankToDB(rankType, rc, rCtx, rankKey, date, guildId); err != nil {
				log.Error("guild in rank to db err", zap.Error(err))
				continue
			}
		}

		delKeys = append(delKeys, key)
	}

	if len(delKeys) > 0 {
		GuildRankInDel(rc, rCtx, delKeys...)
	}

	return nil
}

func FixUpdateGuildInRank() {
	defer func() {
		if r := recover(); r != nil {
			var errField zap.Field
			if err, ok := r.(error); ok {
				errField = zap.String("err", err.Error())
			} else if err, ok := r.(string); ok {
				errField = zap.String("err", err)
			} else {
				errField = zap.Any("err", r)
			}

			stackBuf := make([]byte, 512)
			stackLen := runtime.Stack(stackBuf, true)

			log.Error("FixUpdateGuildInRank panic", errField, zap.ByteString("stack", stackBuf[:stackLen]))
			return
		}
	}()

	var (
		guildId  uint64
		date     int
		err      error
		rankKey  string
		pattern  string
		keySlice []string
		keys     []string

		dbRank *db_guild_in_rank.GuildInRank

		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()

		delKeys = make([]string, 0)
	)

	for _, t := range pb.RankType_value {
		if t == 0 {
			continue
		}
		rankType := pb.RankType(t)
		switch rankType {
		case pb.RankType_RankDaily:
			pattern = enum.REDIS_KEY_GUILD_IN_RANK_DAILY
		case pb.RankType_RankWeekly:
			pattern = enum.REDIS_KEY_GUILD_IN_RANK_WEEKLY
		default:
			log.Error("get guild in rank by type err", zap.Any("rankType", rankType))
			continue
		}
		pattern = pattern + "*"
		keys, err = rc.Keys(rCtx, pattern).Result()
		if err != nil && !errors.Is(err, NIL) {
			log.Error("get guild in rank keys err", zap.Error(err), zap.String("pattern", pattern))
			continue
		}
		if len(keys) <= 0 {
			continue
		}

		for _, key := range keys {
			keySlice = strings.Split(key, ":")

			guildId, err = strconv.ParseUint(keySlice[1], 10, 64)
			if err != nil {
				log.Error("rank date err", zap.String("key", key), zap.Error(err))
				continue
			}

			switch rankType {
			case pb.RankType_RankDaily:
				if FormatGuildInRankDaily(guildId) == key {
					continue
				}
			case pb.RankType_RankWeekly:
				if FormatGuildInRankWeekly(guildId) == key {
					continue
				}
			}

			date, err = strconv.Atoi(keySlice[2])
			if err != nil {
				log.Error("rank date err", zap.String("key", key), zap.Error(err))
				continue
			}

			switch rankType {
			case pb.RankType_RankDaily:
				rankKey = FormatGuildInRankDaily(guildId, date)
			case pb.RankType_RankWeekly:
				rankKey = FormatGuildInRankWeekly(guildId, date)
			}
			dbRank, err = db_guild_in_rank.GetGuildInRankModel().GetRankInfoByTypeAndDateAndGid(rankType, date, guildId)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				log.Error("get guild in rank info err", zap.Error(err))
				continue
			}
			if dbRank == nil || dbRank.ID == 0 {
				if err = GuildInRankToDB(rankType, rc, rCtx, rankKey, date, guildId); err != nil {
					log.Error("guild in rank to db err", zap.Error(err))
					continue
				}
			}

			delKeys = append(delKeys, key)
		}
	}

	if len(delKeys) > 0 {
		GuildRankInDel(rc, rCtx, delKeys...)
	}
}
