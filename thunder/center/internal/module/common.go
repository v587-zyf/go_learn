package module

import (
	"comm/t_data/redis"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
)

func GetRdbUser(userID uint64, locker *rdb_cluster.Locker) (*redis.User, bool) {
	rdbUser, err := redis.GetUser(userID, locker)
	if err != nil {
		log.Error("get redis user err", zap.Error(err), zap.Uint64("userID", userID))
		//grpc_msg.SendErr2User(userID, errcode.ERR_USER_DATA_NOT_FOUND)
		return nil, false
	}

	return rdbUser, true
}
