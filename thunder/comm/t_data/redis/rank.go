package redis

import (
	"comm/t_data/db"
	"comm/t_data/db/db_rank"
	enum "comm/t_enum"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
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

const (
	NIL = redis.Nil
)

func AddGold(userID uint64, gold float64, lv int) error {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	if err := rc.ZIncrBy(rCtx, FormatRankDaily(lv), gold, strconv.Itoa(int(userID))).Err(); err != nil {
		log.Error("add daily rank err", zap.Error(err))
		return err
	}

	if err := rc.ZIncrBy(rCtx, FormatRankWeekly(lv), gold, strconv.Itoa(int(userID))).Err(); err != nil {
		log.Error("add weekly rank err", zap.Error(err))
		return err
	}

	return nil
}

func RankDelOldData(userID uint64, lv int) error {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	if _, err := rc.ZRem(rCtx, FormatRankDaily(lv), strconv.Itoa(int(userID))).Result(); err != nil {
		log.Error("rem daily old rank err", zap.Uint64("userID", userID), zap.Int("lv", lv), zap.Error(err))
		return err
	}

	if _, err := rc.ZRem(rCtx, FormatRankWeekly(lv), strconv.Itoa(int(userID))).Result(); err != nil {
		log.Error("rem weekly old rank err", zap.Uint64("userID", userID), zap.Int("lv", lv), zap.Error(err))
		return err
	}

	return nil
}

func GetRankByType(rankType pb.RankType, lv int) ([]redis.Z, error) {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	var rankKey string
	switch rankType {
	case pb.RankType_RankDaily:
		rankKey = FormatRankDaily(lv)
	case pb.RankType_RankWeekly:
		rankKey = FormatRankWeekly(lv)
	default:
		log.Error("get rank by type err", zap.Any("rankType", rankType), zap.Int("lv", lv))
		return nil, errcode.ERR_PARAM
	}

	return rc.ZRevRangeWithScores(rCtx, rankKey, 0, 499).Result()
}

func GetUserRankingByType(rankType pb.RankType, userID uint64, lv int) (int64, error) {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	var rankKey string
	switch rankType {
	case pb.RankType_RankDaily:
		rankKey = FormatRankDaily(lv)
	case pb.RankType_RankWeekly:
		rankKey = FormatRankWeekly(lv)
	default:
		log.Error("get rank by type err", zap.Any("rankType", rankType))
		return -1, errcode.ERR_PARAM
	}

	member := strconv.Itoa(int(userID))
	return rc.ZRevRank(rCtx, rankKey, member).Result()
}

func GetUserRankScoreByType(rankType pb.RankType, userID uint64, lv int) (float64, error) {
	rc := rdb_single.Get()
	rCtx := rdb_single.GetCtx()

	var rankKey string
	switch rankType {
	case pb.RankType_RankDaily:
		rankKey = FormatRankDaily(lv)
	case pb.RankType_RankWeekly:
		rankKey = FormatRankWeekly(lv)
	default:
		log.Error("get rank by type err", zap.Any("rankType", rankType))
		return -1, errcode.ERR_PARAM
	}

	member := strconv.Itoa(int(userID))
	return rc.ZScore(rCtx, rankKey, member).Result()
}

func RankToDB(rankType pb.RankType, rc *redis.Client, rCtx context.Context, rankKey string, date, lv int) error {
	rdbRankDatas, err := rc.ZRevRangeWithScores(rCtx, rankKey, 0, 99).Result()
	if err != nil {
		log.Error("Error fetching members", zap.Error(err))
		return nil
	}
	if len(rdbRankDatas) <= 0 {
		//log.Warn("rank to db data nil", zap.String("rankKey", rankKey))
		return nil
	}

	seqId, err := db.GenRankIdSeq()
	if err != nil {
		log.Error("gen rank id err", zap.Error(err))
		return err
	}
	addDbData := &db_rank.Rank{
		ID:   seqId,
		Type: rankType,
		Date: date,
		Lv:   lv,
	}

	addDbRanks := make([]*db_rank.RankUnit, len(rdbRankDatas))
	for k, v := range rdbRankDatas {
		userID, err := strconv.ParseUint(v.Member.(string), 10, 64)
		if err != nil {
			log.Error("strconv.ParseUint err", zap.Error(err))
			continue
		}
		addDbRanks[k] = &db_rank.RankUnit{
			Ranking: k + 1,
			UserId:  userID,
			Gold:    int(v.Score),
		}
	}
	addDbData.Rank = addDbRanks
	if _, err = db_rank.GetRankModel().Upsert(addDbData); err != nil {
		log.Error("upsert rank err", zap.Error(err))
		return err
	}
	return nil
}

func RankDel(rc *redis.Client, rCtx context.Context, rankKey ...string) {
	if _, err := rc.Del(rCtx, rankKey...).Result(); err != nil {
		log.Error("del rank err", zap.Error(err))
	}
}

func UpdateRankByRankType(rankType pb.RankType) error {
	var (
		rankKey string
		date    int

		rc   = rdb_single.Get()
		rCtx = rdb_single.GetCtx()

		lvs     = t_tdb.GetLvs()
		delKeys = make([]string, 0)
	)
	for _, lv := range lvs {
		switch rankType {
		case pb.RankType_RankDaily:
			rankKey = FormatRankDaily(lv)
			date = utils.GetYearMonthDay(utils.GetNowUTC())
		case pb.RankType_RankWeekly:
			rankKey = FormatRankWeekly(lv)
			date = utils.GetYearWeek(utils.GetNowUTC())
		default:
			log.Error("get rank by type err", zap.Any("rankType", rankType))
			return fmt.Errorf("get rank by type err")
		}

		//if err := rc.ZRemRangeByRank(rdb.GetCtx(), rankKey, 101, -1).Err(); err != nil {
		//	log.Error("Error removing daily rank members:", zap.Error(err))
		//	return err
		//}

		if err := RankToDB(rankType, rc, rCtx, rankKey, date, lv); err != nil {
			log.Error("rank to db err", zap.Error(err))
			continue
		}

		delKeys = append(delKeys, rankKey)
	}

	if len(delKeys) > 0 {
		RankDel(rc, rCtx, delKeys...)
	}

	return nil
}

func FixUpdateRank() {
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

			log.Error("FixUpdateRank panic", errField, zap.ByteString("stack", stackBuf[:stackLen]))
			return
		}
	}()

	var (
		lv       int
		date     int
		err      error
		rankKey  string
		pattern  string
		keySlice []string
		keys     []string

		dbRank *db_rank.Rank

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
			pattern = enum.REDIS_KEY_RANK_DAILY
		case pb.RankType_RankWeekly:
			pattern = enum.REDIS_KEY_RANK_WEEKLY
		default:
			log.Error("get rank by type err", zap.Any("rankType", rankType))
			continue
		}
		pattern = pattern + "*"
		keys, err = rc.Keys(rCtx, pattern).Result()
		if err != nil && !errors.Is(err, NIL) {
			log.Error("get rank keys err", zap.Error(err), zap.String("pattern", pattern))
			continue
		}
		if len(keys) <= 0 {
			continue
		}

		for _, key := range keys {
			keySlice = strings.Split(key, ":")

			lv, err = strconv.Atoi(keySlice[1])
			if err != nil {
				log.Error("rank date err", zap.String("key", key), zap.Error(err))
				continue
			}
			switch rankType {
			case pb.RankType_RankDaily:
				if FormatRankDaily(lv) == key {
					continue
				}
			case pb.RankType_RankWeekly:
				if FormatRankWeekly(lv) == key {
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
				rankKey = FormatRankDaily(lv, date)
			case pb.RankType_RankWeekly:
				rankKey = FormatRankWeekly(lv, date)
			}
			dbRank, err = db_rank.GetRankModel().GetRankInfoByTypeAndLvAndDate(rankType, date, lv)
			if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
				log.Error("get rank info err", zap.Error(err))
				continue
			}
			if dbRank == nil || dbRank.ID == 0 {
				if err = RankToDB(rankType, rc, rCtx, rankKey, date, lv); err != nil {
					log.Error("rank to db err", zap.Error(err))
					continue
				}
			}

			delKeys = append(delKeys, key)
		}
	}

	if len(delKeys) > 0 {
		RankDel(rc, rCtx, delKeys...)
	}
}
