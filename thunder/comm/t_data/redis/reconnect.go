package redis

import (
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
	"time"
)

func SetReconnect(UserID uint64) (err error) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	rc.Set(rCtx, FormatUserReconnect(UserID), time.Now().Unix(), time.Minute*2)

	return nil
}

func CanReconnect(UserID uint64) bool {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	key := FormatUserReconnect(UserID)
	reconnect, err := rc.Get(rCtx, key).Int64()
	if err != nil {
		log.Error("get Reconnect key err", zap.Error(err))
		return false
	}

	now := time.Now().Unix()
	if now-reconnect > 120 {
		log.Error("Reconnect TimeOut", zap.Uint64("userID", UserID),
			zap.Int64("reconnect", reconnect), zap.Int64("now", now))
		rc.Del(rCtx, key)
		return false
	}

	return true
}

func DelReconnect(UserID uint64) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	rc.Del(rCtx, FormatUserReconnect(UserID))
}
