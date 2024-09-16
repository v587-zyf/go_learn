package db_user

import (
	"comm/t_data/db"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *UserMongoModel) GetUserById(id uint64) (*User, error) {
	var err error

	data := new(User)
	filter := bson.M{"_id": id}
	err = GetUserCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *UserMongoModel) GetUserByChannelInfo(channelInfo *db.AccountChannelInfo) (*User, error) {
	filter := db.MakeAccountFilter(channelInfo)

	data := new(User)
	err := GetUserCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *UserMongoModel) GetUserCount(filter interface{}) (int64, error) {
	return GetUserCol().Find(context.Background(), filter).Count()
}
