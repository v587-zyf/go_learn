package module

import (
	"comm/comm"
	"comm/t_data/db"
	"comm/t_data/redis"
	"comm/t_tdb"
	"fmt"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/rdb/rdb_cluster"
	"go.uber.org/zap"
	"math"
)

func MadeMap(lv int) (newMap *db.Map, err error) {
	userCfg := t_tdb.GetUserUserCfg(lv)
	if userCfg == nil {
		err = fmt.Errorf("user cfg nil. lv:%v", lv)
		return
	}
	mapCfg := t_tdb.GetMapMapCfg(userCfg.Map_id)
	if mapCfg == nil {
		err = fmt.Errorf("map cfg nil. map_id:%v", userCfg.Map_id)
		return
	}

	tempMap := db.NewMap(0, mapCfg.Long, mapCfg.Long, 0, mapCfg.Width, mapCfg.Width)
	// 随机一个出生点
	tempMap.ID = mapCfg.Id
	birthGid := tempMap.RandBirthPos()
	tempMap.SetBirthGid(birthGid)
	tempMap.SetNowGid(birthGid)
	birthGrid := tempMap.GetGridByGid(birthGid)
	birthGrid.RemoveObj(iface.Wall)
	birthGrid.AddObj(iface.Player)
	// 出生点周围九个安全区
	for _, grid := range tempMap.GetSurroundGridByGid(birthGid) {
		grid.SetSafe(true)
	}

	// 生成(长*宽*15.6%)个地雷，四舍五入
	thunderCnt := int(math.Ceil(float64(mapCfg.Long*mapCfg.Width) * float64(t_tdb.Conf().Thunder_rate) / 10000.0))
	//log.Debug("---", zap.Int("long", mapCfg.Long), zap.Int("width", mapCfg.Width),
	//	zap.Float64("cnt", float64(mapCfg.Long*mapCfg.Width)),
	//	zap.Float64("rate", float64(t_tdb.Conf().Thunder_rate)/10000.0),
	//	zap.Int("thunderCnt", thunderCnt))
	for i := 0; i < thunderCnt; i++ {
		thunderGid := tempMap.RandThunderPos()
		tempMap.AddObjToGrid(thunderGid, iface.Thunder)
	}

	// 先生成(长*宽*5%)个宝箱，四舍五
	treasureCnt := int(math.Ceil(float64(mapCfg.Long*mapCfg.Width) * float64(t_tdb.Conf().Treasure_rate) / 10000.0))
	//log.Debug("---", zap.Int("long", mapCfg.Long), zap.Int("width", mapCfg.Width),
	//	zap.Float64("cnt", float64(mapCfg.Long*mapCfg.Width)),
	//	zap.Float64("rate", float64(t_tdb.Conf().Treasure_rate)/10000.0),
	//	zap.Int("treasureCnt", treasureCnt))
	for i := 0; i < treasureCnt; i++ {
		treasureGid := tempMap.RandTreasurePos()
		tempMap.AddObjToGrid(treasureGid, iface.Treasure)
	}
	tempMap.SetTreasure(treasureCnt)

	// 计算所有地雷数量
	// 计算宝箱 当前玩家等级对应的挖墙消耗体力量*随机到的倍率
	grids := tempMap.Grids
	for gid, point := range grids {
		n := 0
		for _, g := range tempMap.GetSurroundGridByGid(gid) {
			if g.HasObj(iface.Thunder) {
				n++
			}
		}
		point.SetThunder(n)

		// 宝箱金币
		if point.HasObj(iface.Treasure) {
			point.SetTreasure(userCfg.Dig_strength * t_tdb.RandTreasureFold())
		}
	}

	newMap = tempMap

	return
}

func GetRdbUser(userID uint64, locker *rdb_cluster.Locker) (*redis.User, bool) {
	rdbUser, err := redis.GetUser(userID, locker)
	if err != nil {
		log.Error("get redis user err", zap.Error(err), zap.Uint64("userID", userID))
		comm.SendErr2User(userID, errcode.ERR_USER_DATA_NOT_FOUND)
		return nil, false
	}

	return rdbUser, true
}
