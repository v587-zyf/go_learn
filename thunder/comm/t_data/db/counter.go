package db

import (
	"context"
	"errors"
	"github.com/v587-zyf/gc/db/mongo"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type Counter struct {
	Name string `bson:"_id"`
	Seq  uint64
}

const (
	COUNTER_USER_ID          = "user_id"
	COUNTER_RANK_ID          = "rank_id"
	COUNTER_GUILD_ID         = "guild_id"
	COUNTER_GUILD_RANK_ID    = "guild_rank_id"
	COUNTER_GUILD_IN_RANK_ID = "guild_in_rank_id"
)

func GenUserIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_USER_ID, 10000000)
}

func GenRankIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_RANK_ID, 0)
}

func GenGuildIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_GUILD_ID, 0)
}

func GenGuildRankIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_GUILD_RANK_ID, 0)
}

func GenGuildInRankIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_GUILD_IN_RANK_ID, 0)
}

func genSeqByName(name string, initValue uint64) (uint64, error) {
	col := mongo.DB(DB_MOKOKO_THUNDER).Collection(COL_COUNTER)

	ret := Counter{}
	err := col.Find(context.Background(), bson.M{"_id": name}).Apply(qmgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		Upsert:    false,
		ReturnNew: true,
	}, &ret)

	if errors.Is(err, qmgo.ErrNoSuchDocuments) {
		_, err = col.InsertOne(context.Background(), &Counter{Name: name, Seq: initValue})
		if err != nil && !qmgo.IsDup(err) {
			return 0, err
		}

		err = col.Find(context.Background(), bson.M{"_id": name}).Apply(qmgo.Change{
			Update:    bson.M{"$inc": bson.M{"seq": 1}},
			Upsert:    false,
			ReturnNew: true,
		}, &ret)
	}

	return ret.Seq, err
}
