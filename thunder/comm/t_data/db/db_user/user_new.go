package db_user

import (
	"comm/t_data/db"
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func (m *UserMongoModel) NewUserUnique(accountInfo *db.AccountChannelInfo, data *User) error {
	filter := db.MakeAccountFilter(accountInfo)

	count, err := GetUserCol().Find(context.Background(), filter).Count()
	if err != nil {
		log.Error("user get account err", zap.Reflect("filter", filter), zap.String("err", err.Error()))
		return err
	}

	if count != 0 {
		return fmt.Errorf("accountRegisterAlready")
	}

	if _, err = GetUserCol().InsertOne(context.Background(), data); err != nil {
		if qmgo.IsDup(err) {
			log.Error("user insert duplicate err", zap.Uint64("id", data.ID), zap.Error(err))
		} else {
			log.Error("user get account err", zap.Reflect("filter", filter), zap.Error(err))
		}

		return err
	}

	return nil
}

func (m *UserMongoModel) InsertOne(data *User) (uint64, error) {
	result, err := GetUserCol().InsertOne(context.Background(), data)
	if err != nil {
		log.Error("user insertOne err", zap.Error(err))
		return 0, err
	}
	insertID := uint64(result.InsertedID.(int64))
	return insertID, nil
}
