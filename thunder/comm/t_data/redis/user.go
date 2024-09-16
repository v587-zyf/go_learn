package redis

import (
	"comm/t_data/db"
	"comm/t_data/db/db_user"
	pb "comm/t_proto/out/client"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
	"time"
)

type User struct {
	ID       uint64       `bson:"_id" json:"_id"`           // 用户id
	Basic    *db.Basic    `bson:"basic" json:"basic"`       // 基础信息
	Telegram *db.Telegram `bson:"telegram" json:"telegram"` // tg
	Map      *db.Map      `bson:"map" json:"map"`           // 地图
	Card     *db.Card     `bson:"card" json:"card"`         // 卡牌
	Shop     *db.Shop     `bson:"shop" json:"shop"`         // 商城
	Hasten   *db.Hasten   `bson:"hasten" json:"hasten"`     // 加速
	Invite   *db.Invite   `bson:"invite" json:"invite"`     // 邀请
	//RedPoint *db.RedPoint `bson:"red_point" json:"red_point"` // 红点
}

func (u *User) ToJson() ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(u)
}

func (u *User) LoadJson(json []byte) error {
	if len(json) == 0 {
		return nil
	}
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(json, u)
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
	//if u.RedPoint == nil {
	//	u.RedPoint = db.NewRedPoint()
	//}
}

func (u *User) LessStrength(lessNum int) {
	if u.Basic.Strength >= lessNum {
		u.Basic.Strength -= lessNum
	} else if u.Basic.ExtraStrength >= lessNum {
		u.Basic.ExtraStrength -= lessNum
	} else if (u.Basic.Strength + u.Basic.ExtraStrength) >= lessNum {
		needStrength := lessNum - u.Basic.Strength
		u.Basic.Strength = 0
		u.Basic.ExtraStrength -= needStrength
	}
}

func (u *User) AddGold(d float64) {
	u.Basic.Gold += d

	if err := AddGold(u.ID, d, u.Basic.Lv); err != nil {
		log.Error("redis add gold err", zap.Uint64("userID", u.ID), zap.Error(err))
	}

	if u.Basic.GuildID != 0 {
		if err := AddGuildGold(u.ID, u.Basic.GuildID, d); err != nil {
			log.Error("redis add guild gold err", zap.Uint64("userID", u.ID), zap.Uint64("guildID", u.Basic.GuildID))
		}
	}
}
func (u *User) LessGold(d float64) {
	u.Basic.Gold -= d
}

func (u *User) AddDiamond(d int) {
	u.Basic.Diamond += d
}
func (u *User) LessDiamond(d int) {
	u.Basic.Diamond -= d
}

func LockUser(userID uint64, SID int64) *rdb_cluster.Locker {
	return rdb_cluster.GetLocker(fmt.Sprint(SID, utils.GUID()), FormatUserID(userID))
}
func GetUser(userID uint64, locker *rdb_cluster.Locker) (*User, error) {
	rc := locker.Get()
	rCtx := locker.GetCtx()

	u := &User{}

	bytes, err := rc.HGet(rCtx, FormatKeyUserID(userID), FormatUserID(userID)).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		return nil, err
	}
	if len(bytes) != 0 {
		// have data
		if err = u.LoadJson([]byte(bytes)); err != nil {
			return nil, err
		}
	} else {
		// create data
		if userDbInfo, err := db_user.GetUserModel().GetUserById(userID); err != nil || userDbInfo.ID == 0 {
			log.Error("get userInfo err", zap.Error(err))
			return nil, err
		} else {
			u, err = UpdateUserByDB(userDbInfo, locker)
			if err != nil {
				log.Error("update userInfo err", zap.Error(err))
				return nil, err
			}
		}
	}

	return u, nil
}
func SetUser(rdbUser *User, locker *rdb_cluster.Locker) error {
	rc := locker.Get()
	rCtx := locker.GetCtx()

	userID := rdbUser.ID
	bytes, err := rdbUser.ToJson()
	if err != nil {
		return err
	}

	if err = rc.HSet(rCtx, FormatKeyUserID(userID), FormatUserID(userID), bytes).Err(); err != nil {
		return err
	}
	rc.Persist(rCtx, FormatKeyUserID(userID))

	z := redis.Z{Score: float64(time.Now().Unix()), Member: userID}
	if err = rc.ZAddNX(rCtx, FormatUserDumpKey(), z).Err(); err != nil {
		//log.Error("ZAddNX userDump err", zap.Error(err), zap.Uint64("userID", userID))
		return err
	}

	return nil
}

func BuildPbMap(uMap *db.Map) *pb.Maps {
	birthPoint := uMap.GetGridByGid(uMap.BirthId)
	nowPoint := uMap.GetGridByGid(uMap.NowId)
	pbMap := &pb.Maps{
		BirthX:   int32(birthPoint.X),
		BirthY:   int32(birthPoint.Y),
		Treasure: int32(uMap.GetTreasure()),
		NowX:     int32(nowPoint.X),
		NowY:     int32(nowPoint.Y),
		Grids:    make([]*pb.MapUnit, len(uMap.Grids)),
	}

	i := 0
	var treasureN int32 = 0
	for _, point := range uMap.Grids {
		if point.HasObj(iface.Treasure) {
			treasureN++
		}
		pbMap.Grids[i] = point.ToPb()
		i++
	}
	pbMap.ReTreasure = treasureN

	return pbMap
}
