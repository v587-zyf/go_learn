package handle

import (
	"comm/comm"
	"comm/t_data/redis"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func GetRdbUser(userID uint64) (*redis.User, bool) {
	locker := redis.LockUser(userID, GetHandleOps().SID)
	if locker == nil {
		log.Error("get locker err")
		return nil, false
	}
	defer locker.Unlock()

	rdbUser, err := redis.GetUser(userID, locker)
	if err != nil {
		log.Error("get redis user err", zap.Error(err), zap.Uint64("userID", userID))
		comm.SendErr2User(userID, errcode.ERR_USER_DATA_NOT_FOUND)
		return nil, false
	}

	return rdbUser, true
}
