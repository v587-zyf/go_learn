package db_guild_rank

import (
	"comm/t_data/db"
	pb "comm/t_proto/out/client"
	"context"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/db/mongo"
	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type GuildRankUnit struct {
	Ranking int    `bson:"ranking"`  // 名次
	GuildId uint64 `bson:"guild_id"` // 公会id
	Gold    int    `bson:"gold"`     // 金币
}
type GuildRank struct {
	ID   uint64           `bson:"_id"`  // 自增id
	Type pb.RankType      `bson:"type"` // 类型
	Date int              `bson:"date"` // 日期
	Rank []*GuildRankUnit `bson:"rank"` // 排名信息
}
type RankMongoModel struct{}

var (
	RankModel = &RankMongoModel{}
)

func GetGuildRankModel() *RankMongoModel {
	return RankModel
}

func GetRankCol() *qmgo.Collection {
	return mongo.DB(db.DB_MOKOKO_THUNDER).Collection(db.COL_GUILD_RANK)
}

func (m *RankMongoModel) InsertOne(data *GuildRank) (uint64, error) {
	result, err := GetRankCol().InsertOne(context.Background(), data)
	if err != nil {
		log.Error("Guild Rank insertOne err", zap.Error(err))
		return 0, err
	}
	insertID := uint64(result.InsertedID.(int64))
	return insertID, nil
}

func (m *RankMongoModel) GetRankInfoById(id uint64) (*GuildRank, error) {
	var data *GuildRank
	var err error
	filter := bson.M{"_id": id}
	err = GetRankCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *RankMongoModel) GetRankInfoByTypeAndDate(t pb.RankType, d int) (*GuildRank, error) {
	var data *GuildRank
	var err error
	filter := bson.M{"type": t, "date": d}
	err = GetRankCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *RankMongoModel) Upsert(data *GuildRank) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	res, err := GetRankCol().Upsert(context.Background(), filter, data)
	if err != nil {
		log.Error("Guild Rank upsert err", zap.Reflect("data", data), zap.String("err", err.Error()))
		return nil, err
	}
	return res, nil
}
