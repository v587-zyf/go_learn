package db_rank

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

type RankUnit struct {
	Ranking int    `bson:"ranking"` // 名次
	UserId  uint64 `bson:"user_id"` // 用户id
	Gold    int    `bson:"gold"`    // 金币
}
type Rank struct {
	ID   uint64      `bson:"_id"`  // 自增id
	Type pb.RankType `bson:"type"` // 类型
	Date int         `bson:"date"` // 日期
	Lv   int         `bson:"lv"`   // 等级
	Rank []*RankUnit `bson:"rank"` // 排名信息
}
type RankMongoModel struct{}

var (
	RankModel = &RankMongoModel{}
)

func GetRankModel() *RankMongoModel {
	return RankModel
}

func GetRankCol() *qmgo.Collection {
	return mongo.DB(db.DB_MOKOKO_THUNDER).Collection(db.COL_RANK)
}

func (m *RankMongoModel) InsertOne(data *Rank) (uint64, error) {
	result, err := GetRankCol().InsertOne(context.Background(), data)
	if err != nil {
		log.Error("Rank insertOne err", zap.Error(err))
		return 0, err
	}
	insertID := uint64(result.InsertedID.(int64))
	return insertID, nil
}

func (m *RankMongoModel) GetRankInfoById(id uint64) (*Rank, error) {
	var data *Rank
	var err error
	filter := bson.M{"_id": id}
	err = GetRankCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *RankMongoModel) GetRankInfoByTypeAndLvAndDate(t pb.RankType, d, lv int) (*Rank, error) {
	var data *Rank
	var err error
	filter := bson.M{"type": t, "date": d, "lv": lv}
	err = GetRankCol().Find(context.Background(), filter).One(&data)
	return data, err
}

func (m *RankMongoModel) Upsert(data *Rank) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	res, err := GetRankCol().Upsert(context.Background(), filter, data)
	if err != nil {
		log.Error("Rank upsert err", zap.Reflect("data", data), zap.String("err", err.Error()))
		return nil, err
	}
	return res, nil
}
