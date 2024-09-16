package db_guild

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

func (m *GuildMongoModel) NewGuildUnique(data *Guild) error {
	filter := bson.M{"chat_id": data.ChatID}
	count, err := GetGuildCol().Find(context.Background(), filter).Count()
	if err != nil {
		log.Error("guild get chatId err", zap.Reflect("filter", filter), zap.String("err", err.Error()))
		return err
	}

	if count != 0 {
		return fmt.Errorf("GuildRegisterAlready")
	}

	if _, err = GetGuildCol().InsertOne(context.Background(), data); err != nil {
		if qmgo.IsDup(err) {
			log.Error("guild insert duplicate err", zap.Uint64("id", data.ID), zap.Error(err))
		} else {
			log.Error("guild get account err", zap.Reflect("filter", filter), zap.Error(err))
		}

		return err
	}

	return nil
}

func (m *GuildMongoModel) InsertOne(data *Guild) (uint64, error) {
	result, err := GetGuildCol().InsertOne(context.Background(), data)
	if err != nil {
		log.Error("guild insertOne err", zap.Error(err))
		return 0, err
	}
	insertID := uint64(result.InsertedID.(int64))
	return insertID, nil
}
