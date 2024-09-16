package redis

import (
	"comm/t_enum"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func DbDump(SID int64) error {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	ctx, cancel := context.WithTimeout(rCtx, 5*time.Second)
	defer cancel()

	t := time.Now().Add(-30 * time.Second)

	userIDs, err := rc.ZRangeByScore(ctx, FormatUserDumpKey(), &redis.ZRangeBy{
		Min:   "0",
		Max:   strconv.FormatInt(t.Unix(), 10),
		Count: 100,
	}).Result()
	if err != nil {
		log.Error("doDump: redis err", zap.String("err", err.Error()))
	}
	if len(userIDs) <= 0 {
		return nil
	}

	remUIDs := make([]any, 0)
	for _, v := range userIDs {
		userID := utils.StrToUInt64(v)
		if userID == 0 {
			log.Error("invalid userID", zap.String("userID", v))
			continue
		}

		// update db
		if err = UpdateUserDbByRedis(userID, SID); err != nil {
			log.Error("db upsert err", zap.Uint64("userID", userID), zap.Error(err))
			continue
		}

		// set expire
		rc.Expire(rCtx, FormatKeyUserID(userID), enum.DUMP_USER_EXPIRE_TIME)

		remUIDs = append(remUIDs, v)

		//log.Debug("user dump wealth to mongo", zap.Uint64("userID", userID))
	}

	if err = rc.ZRem(rCtx, FormatUserDumpKey(), remUIDs...).Err(); err != nil {
		log.Error("redis ZRem dbDump err", zap.String("err", err.Error()))
		return err
	}

	return nil
}
