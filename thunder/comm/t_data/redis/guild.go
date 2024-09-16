package redis

import (
	"comm/t_data/db"
	"comm/t_data/db/db_guild"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"github.com/v587-zyf/gc/utils"
	"go.uber.org/zap"
)

type Guild struct {
	ID        uint64           `bson:"_id" json:"_id"`               // 自增id
	CreateUID uint64           `bson:"create_uid" json:"create_uid"` // 创建者id
	ChatID    int64            `bson:"chat_id" json:"chat_id"`       // 聊天id
	Title     string           `bson:"title" json:"title"`           // 公会名称
	Members   *db.GuildMembers `bson:"members" json:"members"`       // 成员
}

func (g *Guild) ToJson() ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(g)
}

func (g *Guild) LoadJson(json []byte) error {
	if len(json) == 0 {
		return nil
	}
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(json, g)
}

func (g *Guild) AddMember(UID uint64) {
	g.Members.Add(UID)
}
func (g *Guild) DelMember(UID uint64) {
	g.Members.Del(UID)
}

func LockGuild(guildID uint64, SID int64) *rdb_cluster.Locker {
	return rdb_cluster.GetLocker(fmt.Sprint(SID, utils.GUID()), FormatGuildID(guildID))
}
func GetGuild(guildID uint64, locker *rdb_cluster.Locker) (*Guild, error) {
	rc := locker.Get()
	rCtx := locker.GetCtx()

	g := &Guild{}
	bytes, err := rc.HGet(rCtx, FormatGuild(), FormatGuildID(guildID)).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		return nil, err
	}
	if len(bytes) != 0 {
		if err = g.LoadJson([]byte(bytes)); err != nil {
			return nil, err
		}
	} else {
		if DbGuildInfo, err := db_guild.GetGuildModel().GetGuildById(guildID); err != nil || DbGuildInfo.ID == 0 {
			log.Error("get guildInfo err", zap.Error(err))
			return nil, err
		} else {
			g.ID = DbGuildInfo.ID
			g.CreateUID = DbGuildInfo.CreateUID
			g.ChatID = DbGuildInfo.ChatID
			g.Title = DbGuildInfo.Title
			g.Members = DbGuildInfo.Members

			if err = SetGuild(g, locker); err != nil {
				log.Error("redis write guild err", zap.Error(err))
				return nil, err
			}
		}
	}

	return g, nil
}

func SetGuild(data *Guild, locker *rdb_cluster.Locker) (err error) {
	rc := locker.Get()
	rCtx := locker.GetCtx()

	bytes, err := data.ToJson()
	if err != nil {
		return err
	}

	if _, err = rc.HSet(rCtx, FormatGuild(), FormatGuildID(data.ID), bytes).Result(); err != nil {
		log.Error("Add guild err", zap.Error(err))
		return
	}

	return nil
}

func GetAllGuild() ([]*Guild, error) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	result, err := rc.HMGet(rCtx, FormatGuild()).Result()
	if err != nil {
		log.Error("ger all guild err", zap.Error(err))
		return nil, err
	}

	var guilds []*Guild
	for _, v := range result {
		if v == nil {
			continue
		}
		guild := &Guild{}
		if err = guild.LoadJson(v.([]byte)); err != nil {
			log.Error("load guild err", zap.Error(err), zap.Any("v", v))
			continue
		}
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func GetGuildByIds(guildID ...string) ([]*Guild, error) {
	rc := rdb_cluster.Get()
	rCtx := rdb_cluster.GetCtx()

	result, err := rc.HMGet(rCtx, FormatGuild(), guildID...).Result()
	if err != nil {
		log.Error("ger all guild err", zap.Error(err))
		return nil, err
	}

	var guilds []*Guild
	for _, v := range result {
		if v == nil {
			continue
		}
		guild := &Guild{}
		if err = guild.LoadJson(v.([]byte)); err != nil {
			log.Error("load guild err", zap.Error(err), zap.Any("v", v))
			continue
		}
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

func AddGuildGold(userID, guildID uint64, gold float64) (err error) {
	if err = AddGuildInRankGold(userID, guildID, gold); err != nil {
		return
	}

	if err = AddGuildRankGold(guildID, gold); err != nil {
		return
	}

	return
}
