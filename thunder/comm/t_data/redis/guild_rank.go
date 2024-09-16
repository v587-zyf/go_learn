package redis

import (
	"comm/t_data/db"
	"comm/t_data/db/db_guild_rank"
	enum "comm/t_enum"
	pb "comm/t_proto/out/client"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_single"
	"github.com/v587-zyf/gc/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"runtime"
	"strconv"
	"strings"
)

func AddGuildRankGold(guildID uint64, gold float64) error {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	if err := rc.ZIncrBy(rCtx, FormatGuildRankDaily(), gold, strconv.Itoa(int(guildID))).Err(); err != nil {
		log.Error("add guild in daily rank err", zap.Error(err))
		return err
	}

	if err := rc.ZIncrBy(rCtx, FormatGuildRankWeekly(), gold, strconv.Itoa(int(guildID))).Err(); err != nil {
		log.Error("add guild in weekly rank err", zap.Error(err))
		return err
	}

	return nil
}

func GetGuildRankByType(rankType pb.RankType) ([]redis.Z, error) {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	var rankKey string
	switch rankType {
	case pb.RankType_RankDaily:
		rankKey = FormatGuildRankDaily()
	case pb.RankType_RankWeekly:
		rankKey = FormatGuildRankWeekly()
	default:
		log.Error("get rank by type err", zap.Any("rankType", rankType))
		return nil, errcode.ERR_PARAM
	}

	return rc.ZRevRangeWithScores(rCtx, rankKey, 0, 499).Result()
}

func GuildRankToDB(rankType pb.RankType, rc *redis.Client, rCtx context.Context, rankKey string, date int) error {
	rdbRankDatas, err := rc.ZRevRangeWithScores(rCtx, rankKey, 0, 99).Result()
	if err != nil {
		log.Error("Error fetching guild rank", zap.Error(err))
		return nil
	}
	if len(rdbRankDatas) <= 0 {
		//log.Warn("rank to db data nil", zap.String("rankKey", rankKey))
		return nil
	}

	seqId, err := db.GenGuildRankIdSeq()
	if err != nil {
		log.Error("gen guild rank id err", zap.Error(err))
		return err
	}
	addDbData := &db_guild_rank.GuildRank{
		ID:   seqId,
		Type: rankType,
		Date: date,
	}

	addDbRanks := make([]*db_guild_rank.GuildRankUnit, len(rdbRankDatas))
	for k, v := range rdbRankDatas {
		guildId, err := strconv.ParseUint(v.Member.(string), 10, 64)
		if err != nil {
			log.Error("strconv.ParseUint err", zap.Error(err))
			continue
		}
		addDbRanks[k] = &db_guild_rank.GuildRankUnit{
			Ranking: k + 1,
			GuildId: guildId,
			Gold:    int(v.Score),
		}
	}
	addDbData.Rank = addDbRanks
	if _, err = db_guild_rank.GetGuildRankModel().Upsert(addDbData); err != nil {
		log.Error("upsert guild rank err", zap.Error(err))
		return err
	}
	return nil
}

func GuildRankDel(rc *redis.Client, rCtx context.Context, rankKey ...string) {
	if _, err := rc.Del(rCtx, rankKey...).Result(); err != nil {
		log.Error("del guild rank err", zap.Error(err))
	}
}

func UpdateGuildRankByRankType(rankType pb.RankType) error {
	var (
		rankKey string
		date    int

		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()
	)

	switch rankType {
	case pb.RankType_RankDaily:
		rankKey = FormatGuildRankDaily()
		date = utils.GetYearMonthDay(utils.GetNowUTC())
	case pb.RankType_RankWeekly:
		rankKey = FormatGuildRankWeekly()
		date = utils.GetYearWeek(utils.GetNowUTC())
	default:
		log.Error("get guild rank by type err", zap.Any("rankType", rankType))
		return fmt.Errorf("get guild rank by type err")
	}

	if err := GuildRankToDB(rankType, rc, rCtx, rankKey, date); err != nil {
		log.Error("guild rank to db err", zap.Error(err))
		return err
	}

	GuildRankDel(rc, rCtx, rankKey)

	return nil
}

func FixUpdateGuildRank() {
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

			log.Error("FixUpdateGuildRank panic", errField, zap.ByteString("stack", stackBuf[:stackLen]))
			return
		}
	}()

	var (
		date     int
		err      error
		rankKey  string
		pattern  string
		keySlice []string
		keys     []string

		dbRank *db_guild_rank.GuildRank

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
			pattern = enum.REDIS_KEY_GUILD_RANK_DAILY
		case pb.RankType_RankWeekly:
			pattern = enum.REDIS_KEY_GUILD_RANK_WEEKLY
		default:
			log.Error("get guild rank by type err", zap.Any("rankType", rankType))
			continue
		}
		pattern = pattern + "*"
		keys, err = rc.Keys(rCtx, pattern).Result()
		if err != nil && !errors.Is(err, NIL) {
			log.Error("get guild rank keys err", zap.Error(err), zap.String("pattern", pattern))
			continue
		}
		if len(keys) <= 0 {
			continue
		}

		for _, key := range keys {
			keySlice = strings.Split(key, ":")

			date, err = strconv.Atoi(keySlice[1])
			if err != nil {
				log.Error("rank date err", zap.String("key", key), zap.Error(err))
				continue
			}

			switch rankType {
			case pb.RankType_RankDaily:
				rankKey = FormatGuildRankDaily(date)
			case pb.RankType_RankWeekly:
				rankKey = FormatGuildRankWeekly(date)
			}
			dbRank, err = db_guild_rank.GetGuildRankModel().GetRankInfoByTypeAndDate(rankType, date)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				log.Error("get guild rank info err", zap.Error(err))
				continue
			}
			if dbRank == nil || dbRank.ID == 0 {
				if err = GuildRankToDB(rankType, rc, rCtx, rankKey, date); err != nil {
					log.Error("GUILD rank to db err", zap.Error(err))
					continue
				}
			}

			delKeys = append(delKeys, key)
		}
	}

	if len(delKeys) > 0 {
		GuildRankDel(rc, rCtx, delKeys...)
	}
}
