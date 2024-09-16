package db_guild

import (
	"comm/t_data/db"
	"context"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"time"
)

type Guild struct {
	ID        uint64           `bson:"_id" json:"_id"`               // 自增id
	CreateUID uint64           `bson:"create_uid" json:"create_uid"` // 创建者id
	ChatID    int64            `bson:"chat_id" json:"chat_id"`       // 聊天id
	Title     string           `bson:"title" json:"title"`           // 公会名称
	Members   *db.GuildMembers `bson:"members" json:"members"`       // 成员
	CreateAt  time.Time        `bson:"create_at" json:"create_at"`   // 创建时间
}

type GuildMongoModel struct{}

var (
	GuildModel = &GuildMongoModel{}
)

func GetGuildModel() *GuildMongoModel {
	return GuildModel
}

func GetGuildCol() *qmgo.Collection {
	return mongo.DB(db.DB_MOKOKO_THUNDER).Collection(db.COL_GUILD)
}

func (m *GuildMongoModel) Upsert(data *Guild) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	res, err := GetGuildCol().Upsert(context.Background(), filter, data)
	if err != nil {
		log.Error("guild upsert err", zap.Reflect("data", data), zap.String("err", err.Error()))
		return nil, err
	}
	return res, nil
}
