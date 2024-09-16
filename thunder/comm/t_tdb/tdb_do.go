package t_tdb

import (
	pb "comm/t_proto/out/client"
	"github.com/v587-zyf/gc/utils"
)

func RandTreasureFold() int {
	weightMap := make(map[int]int)
	for id, cfg := range t_tdb.TreasureTreasureCfgs {
		weightMap[id] = cfg.Rate
	}

	return GetTreasureTreasureCfg(utils.RandWeightByMap(weightMap)).Fold
}

func GetCardsByMapID(mapID int) []*CardCardCfg {
	cfgs := make([]*CardCardCfg, 0)
	for i := 1; i <= mapID; i++ {
		cfgs = append(cfgs, t_tdb.CardMap[i]...)
	}
	return cfgs
}
func RandCard(mapID int, gachaNum int) []int {
	ids := make([]int, gachaNum)
	for i := 0; i < gachaNum; i++ {
		weightMap := make(map[int]int)
		for _, cfg := range GetCardsByMapID(mapID) {
			weightMap[cfg.Id] = cfg.Weight
		}
		ids[i] = utils.RandWeightByMap(weightMap)
	}

	return ids
}

func GetCardLv(cardID, lv int) *CardLvCardLvCfg {
	return t_tdb.CardLvMap[cardID][lv]
}
func GetShopLv(t pb.ShopType, lv int) *ShopShopCfg {
	return t_tdb.ShopLvMap[int(t)][lv]
}

func GetLvs() []int {
	return t_tdb.LvSlice
}
