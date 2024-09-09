package mongo

import (
	"context"
	"kernel/db/mongo"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type CWCounter struct {
	Name string `bson:"_id"`
	Seq  uint64
}

const (
	COUNTER_TEST = "testId"
)

// GenBattleTokenSeq 获取用户自增id
func GetTestIdSeq() (uint64, error) {
	return genSeqByName(COUNTER_TEST, 1)
}

func genSeqByName(name string, initValue uint64) (uint64, error) {
	col := mongo.DB(DB_TEST).Collection(COL_COUNTER)

	ret := CWCounter{}

	err := col.Find(context.Background(), bson.M{"_id": name}).Apply(qmgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		Upsert:    false,
		ReturnNew: true,
	}, &ret)

	if err == qmgo.ErrNoSuchDocuments {
		_, err = col.InsertOne(context.Background(), &CWCounter{Name: name, Seq: initValue})
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
