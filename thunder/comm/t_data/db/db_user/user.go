package db_user

import (
	"comm/t_data/db"
	"context"
	"github.com/qiniu/qmgo"
	"github.com/v587-zyf/gc/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type User struct {
	ID       uint64       `bson:"_id"`      // 用户id
	Basic    *db.Basic    `bson:"basic"`    // 基础信息
	Accounts *db.Accounts `bson:"accounts"` // 账号
	Telegram *db.Telegram `bson:"telegram"` // tg
	Map      *db.Map      `bson:"map"`      // 地图
	Card     *db.Card     `bson:"card"`     // 卡牌
	Shop     *db.Shop     `bson:"shop"`     // 商城
	Hasten   *db.Hasten   `bson:"hasten"`   // 加速
	Invite   *db.Invite   `bson:"invite"`   // 邀请
	//RedPoint *db.RedPoint `bson:"red_point"` // 红点
}
type UserMongoModel struct{}

var (
	UserModel = &UserMongoModel{}
)

func GetUserModel() *UserMongoModel {
	return UserModel
}

func GetUserCol() *qmgo.Collection {
	return mongo.DB(db.DB_MOKOKO_THUNDER).Collection(db.COL_USER)
}

func (u *User) Init() {
	if u.Shop == nil || u.Shop.Data == nil {
		u.Shop = db.NewShop()
	}
	if u.Hasten == nil || u.Hasten.Data == nil {
		u.Hasten = db.NewHasten()
	}
	if u.Invite == nil {
		u.Invite = db.NewInvite(0)
	}
}

func (m *UserMongoModel) Upsert(data *User) (*qmgo.UpdateResult, error) {
	filter := bson.M{"_id": data.ID}
	res, err := GetUserCol().Upsert(context.Background(), filter, data)
	if err != nil {
		log.Error("user upsert err", zap.Reflect("data", data), zap.String("err", err.Error()))
		return nil, err
	}
	return res, nil
}
