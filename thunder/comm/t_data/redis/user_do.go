package redis

import (
	"comm/t_data/db/db_user"
	"fmt"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
	"time"
)

func UpdateUserByDB(userDbInfo *db_user.User, locker *rdb_cluster.Locker) (*User, error) {
	u := new(User)

	userID := userDbInfo.ID
	u.ID = userID

	u.Basic = userDbInfo.Basic
	u.Telegram = userDbInfo.Telegram
	u.Map = userDbInfo.Map
	u.Card = userDbInfo.Card
	u.Shop = userDbInfo.Shop
	u.Hasten = userDbInfo.Hasten
	u.Invite = userDbInfo.Invite
	//u.RedPoint = userDbInfo.RedPoint

	if err := SetUser(u, locker); err != nil {
		log.Error("set userInfo err", zap.Error(err))
		return nil, err
	}

	return u, nil
}

func UpdateUserDbByRedis(userID uint64, SID int64) (err error) {
	locker := LockUser(userID, SID)
	if locker == nil {
		err = fmt.Errorf("get locker err")
		return
	}
	defer locker.Unlock()

	userRedisInfo, err := GetUser(userID, locker)
	if err != nil {
		log.Error("GetRedisUser err", zap.Error(err))
		return
	}

	timeNow := time.Now()

	dbUser := db_user.GetUserModel()
	userDbInfo, err := dbUser.GetUserById(userID)
	if err != nil {
		log.Error("get userDbInfo err", zap.Error(err))
		return
	}

	userDbInfo.Basic = userRedisInfo.Basic
	userDbInfo.Basic.LastUpdateAt = timeNow

	userDbInfo.Map = userRedisInfo.Map
	userDbInfo.Card = userRedisInfo.Card
	userDbInfo.Shop = userRedisInfo.Shop
	userDbInfo.Hasten = userRedisInfo.Hasten
	userDbInfo.Invite = userRedisInfo.Invite
	//userDbInfo.RedPoint = userRedisInfo.RedPoint

	if _, err = dbUser.Upsert(userDbInfo); err != nil {
		log.Error("upsert userInfo err", zap.Error(err))
		return
	}

	return
}
