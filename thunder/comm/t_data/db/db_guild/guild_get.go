package db_guild

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *GuildMongoModel) GetGuildById(id uint64) (*Guild, error) {
	var err error

	data := new(Guild)
	filter := bson.M{"_id": id}
	err = GetGuildCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *GuildMongoModel) GetGuildByChatId(chatId int64) (*Guild, error) {
	var err error

	data := new(Guild)
	filter := bson.M{"chat_id": chatId}
	err = GetGuildCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *GuildMongoModel) GetGuildCount(filter interface{}) (int64, error) {
	return GetGuildCol().Find(context.Background(), filter).Count()
}
