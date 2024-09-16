package db

import (
	"comm/t_tdb"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"time"
)

const (
	DEFAULT_USER_LV = 1
)

type Basic struct {
	Lv             int       `bson:"lv" json:"lv,omitempty"`                             // 等级
	Gold           float64   `bson:"gold" json:"gold,omitempty"`                         // 金币
	LastGoldAt     time.Time `bson:"last_gold_at" json:"last_gold_at,omitempty"`         // 最后收益时间
	Diamond        int       `bson:"diamond" json:"diamond,omitempty"`                   // 钻石
	Strength       int       `bson:"strength" json:"strength,omitempty"`                 // 体力
	ExtraStrength  int       `bson:"extra_strength" json:"extra_strength,omitempty"`     // 额外体力
	LastStrengthAt time.Time `bson:"last_strength_at" json:"last_strength_at,omitempty"` // 最后体力恢复时间
	Dead           int       `bson:"dead" json:"dead,omitempty"`                         // 是否死亡
	DeadPointID    int       `bson:"dead_point_id" json:"dead_point_id,omitempty"`       // 死亡地雷格子id
	TreasureCnt    int       `bson:"treasure_cnt" json:"treasure_cnt,omitempty"`         // 总宝箱数
	FreeResetMap   int       `bson:"free_reset_map" json:"free_reset_map,omitempty"`     // 是否可以免费重置地图
	GuildID        uint64    `bson:"guild_id" json:"guild_id,omitempty"`                 // 公会id
	LastUpdateAt   time.Time `bson:"last_update_at" json:"lastUpdate_at,omitempty"`      // 最后修改时间
	CreateAt       time.Time `bson:"create_at" json:"create_at,omitempty"`               // 创建时间
}

func NewBasic() *Basic {
	userCfg := t_tdb.GetUserUserCfg(DEFAULT_USER_LV)
	if userCfg == nil {
		log.Error("user cfg nil", zap.Int("id", DEFAULT_USER_LV))
		return nil
	}
	return &Basic{
		Lv:           DEFAULT_USER_LV,
		Gold:         0,
		Strength:     userCfg.Strength_max,
		FreeResetMap: enums.NO,
		Dead:         enums.NO,
		CreateAt:     time.Now(),
	}
}
