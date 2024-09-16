package module

import (
	"comm/t_data/redis"
	"comm/t_enum"
	"comm/t_errcode"
	pb "comm/t_proto/out/client"
	"comm/t_tdb"
	"fmt"
	"github.com/v587-zyf/gc/enums"
	"github.com/v587-zyf/gc/errcode"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"github.com/v587-zyf/gc/module"
	"go.uber.org/zap"
	"math"
)

type MapMgr struct {
	module.DefModule
}

func NewMapMgr() *MapMgr {
	return &MapMgr{}
}

func (m *MapMgr) Move(userID uint64, req *pb.MoveReq) (ntf iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	if rdbUser.Basic.Dead == enums.YES {
		err = errCode.ERR_DEAD
		return
	}
	// 不能走墙、宝箱上
	point := rdbUser.Map.GetGridByPos(req.X, req.Y)
	if point == nil || point.HasObj(iface.Wall) || point.HasObj(iface.Player) || point.HasObj(iface.Thunder) {
		err = errCode.ERR_MOVE
		return
	}

	rdbUser.Map.GetGridByGid(rdbUser.Map.GetNowGid()).RemoveObj(iface.Player).AddObj(iface.Empty)
	rdbUser.Map.GetGridByPos(req.X, req.Y).AddObj(iface.Player)
	rdbUser.Map.SetNowXY(float32(req.X), float32(req.Y))

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("update redis user err", zap.Error(err), zap.Uint64("userID", rdbUser.ID))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	//fmt.Printf(rdbUser.Map.String())
	msgID = pb.MsgID_Move_AckId
	ntf = &pb.MoveAck{
		X: req.X,
		Y: req.Y,
	}

	return
}

func (m *MapMgr) OpenWall(userID uint64, req *pb.OpenWallReq) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	userCfg := t_tdb.GetUserUserCfg(rdbUser.Basic.Lv)
	if userCfg == nil {
		err = fmt.Errorf("user cfg nil. lv:%v", rdbUser.Basic.Lv)
		return
	}

	if rdbUser.Basic.Dead == enums.YES {
		err = errCode.ERR_DEAD
		return
	}

	if (rdbUser.Basic.Strength + rdbUser.Basic.ExtraStrength) < userCfg.Dig_strength {
		err = errCode.ERR_STRENGTH_NOT_ENOUGH
		return
	}

	point := rdbUser.Map.GetGridByPos(req.X, req.Y)
	if point == nil || !point.HasObj(iface.Wall) {
		err = errCode.ERR_OPEN_WALL
		return
	}

	rdbUser.LessStrength(userCfg.Dig_strength)

	ackGrids := make([]*pb.MapUnit, 0)
	if point.HasObj(iface.Thunder) {
		point.RemoveObj(iface.Wall)
		rdbUser.Basic.Dead = enums.YES
		rdbUser.Basic.DeadPointID = point.ID()

		ackGrids = append(ackGrids, point.ToPb())
	} else {
		remNum := 1
		filterMap := make(map[int]struct{})
		filterMap = rdbUser.Map.GetSurroundNoThunderPids(point, filterMap)
		if len(filterMap) > 0 {
			var pPb *pb.MapUnit
			for k := range filterMap {
				pPb = rdbUser.Map.GetGridByGid(k).RemoveObj(iface.Wall).AddObj(iface.Empty).ToPb()

				ackGrids = append(ackGrids, pPb)
			}
			remNum = len(filterMap)
		} else {
			point.RemoveObj(iface.Wall)
			point.AddObj(iface.Empty)

			ackGrids = append(ackGrids, point.ToPb())
		}

		//log.Debug("---", zap.Float64("addGold", float64(userCfg.Dig_strength*remNum)), zap.Int("remNum", remNum))
		// 获得金币量=消耗体力量*同时挖开的墙面数量
		rdbUser.AddGold(float64(userCfg.Dig_strength * remNum))

		defer func() {
			GetClientModuleMgr().GetModule(enum.G_M_USER).(*UserMgr).UpLv(rdbUser.ID, locker)
		}()
	}
	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("update redis user err", zap.Error(err), zap.Uint64("userID", rdbUser.ID))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	//fmt.Printf(rdbUser.Map.String())

	msgID = pb.MsgID_OpenWall_AckId
	msg = &pb.OpenWallAck{
		Strength:      int32(rdbUser.Basic.Strength),
		Grids:         ackGrids,
		ExtraStrength: int32(rdbUser.Basic.ExtraStrength),
		Dead:          rdbUser.Basic.Dead == enums.YES,
		Gold:          rdbUser.Basic.Gold,
		IsOver:        rdbUser.Map.IsOver(),
	}
	return
}

func (m *MapMgr) GetTreasure(userID uint64, req *pb.GetTreasureReq) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	point := rdbUser.Map.GetGridByPos(req.X, req.Y)
	if point == nil || !point.HasObj(iface.Treasure) {
		err = errCode.ERR_POINT
		return
	}

	msgID = pb.MsgID_GetTreasure_AckId
	msg = &pb.GetTreasureAck{Gold: float64(point.GetTreasure())}

	return
}

func (m *MapMgr) OpenTreasure(userID uint64, req *pb.OpenTreasureReq) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	point := rdbUser.Map.GetGridByPos(req.X, req.Y)
	if point == nil || !point.HasObj(iface.Treasure) {
		err = errCode.ERR_POINT
		return
	}

	point.RemoveObj(iface.Treasure)
	point.AddObj(iface.Empty)

	//rdbUser.Map.SetTreasure(rdbUser.Map.GetTreasure() - 1)
	rdbUser.AddGold(float64(point.GetTreasure()))

	defer func() {
		GetClientModuleMgr().GetModule(enum.G_M_CARD).(*CardMgr).GetCard(rdbUser.Map.ID, rdbUser.ID, locker)
	}()

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("update redis user err", zap.Error(err), zap.Uint64("userID", rdbUser.ID))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_OpenTreasure_AckId
	treasureN := 0
	for _, pit := range rdbUser.Map.Grids {
		if pit.HasObj(iface.Treasure) {
			treasureN++
		}
	}
	msg = &pb.OpenTreasureAck{Gold: rdbUser.Basic.Gold, Grid: point.ToPb(), Treasure: int32(treasureN)}

	return
}

func (m *MapMgr) Revive(userID uint64, req *pb.ReviveReq) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	if rdbUser.Basic.Dead != enums.YES {
		err = errCode.ERR_NO_DEAD
		return
	}
	if (rdbUser.Basic.Strength + rdbUser.Basic.ExtraStrength) < t_tdb.Conf().Revive_strength {
		err = errCode.ERR_STRENGTH_NOT_ENOUGH
		return
	}

	thunderPoint := rdbUser.Map.GetGridByGid(rdbUser.Basic.DeadPointID)

	//birX, birY := rdbUser.Map.GetBirthXY()
	//rdbUser.Map.SetNowXY(birX, birY)
	nowX, nowY := rdbUser.Map.GetNowXY()
	thunderPoint.RemoveObj(iface.Thunder).AddObj(iface.Empty)

	grids := rdbUser.Map.GetSurroundGrids25(rdbUser.Basic.DeadPointID)
	ackGrids := make([]*pb.MapUnit, len(grids))
	num := 0
	for _, grid := range grids {
		n := 0
		for _, g := range rdbUser.Map.GetSurroundGridByGid(grid.ID()) {
			if g.HasObj(iface.Thunder) {
				n++
			}
		}
		grid.SetThunder(n)
		ackGrids[num] = grid.ToPb()
		num++
	}

	rdbUser.LessStrength(t_tdb.Conf().Revive_strength)
	rdbUser.Basic.Dead = enums.NO
	rdbUser.Basic.DeadPointID = 0

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("update redis user err", zap.Error(err), zap.Uint64("userID", rdbUser.ID))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	msgID = pb.MsgID_Revive_AckId
	msg = &pb.ReviveAck{
		Dead:          rdbUser.Basic.Dead == enums.YES,
		Strength:      int32(rdbUser.Basic.Strength),
		Grid:          thunderPoint.ToPb(),
		NowX:          int32(nowX),
		NowY:          int32(nowY),
		Grids:         ackGrids,
		ExtraStrength: int32(rdbUser.Basic.ExtraStrength),
	}

	return
}

func (m *MapMgr) ResetMap(userID uint64, req *pb.ResetMapReq) (msg iface.IProtoMessage, msgID int32, err error) {
	locker := redis.LockUser(userID, GetClientModuleMgrOptions().SID)
	if locker == nil {
		err = errcode.ERR_USER_DATA_INVALID
		return
	}
	defer locker.Unlock()

	rdbUser, ret := GetRdbUser(userID, locker)
	if !ret {
		err = errcode.ERR_USER_DATA_NOT_FOUND
		return
	}

	userCfg := t_tdb.GetUserUserCfg(rdbUser.Basic.Lv)
	if userCfg == nil {
		err = fmt.Errorf("user cfg nil. lv:%v", rdbUser.Basic.Lv)
		return
	}
	mapCfg := t_tdb.GetMapMapCfg(userCfg.Map_id)
	if mapCfg == nil {
		err = fmt.Errorf("map cfg nil. map_id:%v", userCfg.Map_id)
		return
	}

	// free
	if rdbUser.Basic.FreeResetMap == enums.YES {
		rdbUser.Basic.FreeResetMap = enums.NO
	} else {
		// no free
		if rdbUser.Basic.TreasureCnt != 0 {
			// 消耗量=地图等级对应的重置地图消耗值*当前地图剩余未开启宝箱数量/地图宝箱总量，四舍五入（也就是说宝箱全开完后重置地图就没有消耗了）
			rmStrength := int(math.Ceil(float64(mapCfg.Reset_strength * rdbUser.Map.Treasure / rdbUser.Basic.TreasureCnt)))
			if (rdbUser.Basic.Strength + rdbUser.Basic.ExtraStrength) < rmStrength {
				err = errCode.ERR_STRENGTH_NOT_ENOUGH
				return
			}

			rdbUser.LessStrength(rmStrength)
		}
	}

	rdbUser.Map, err = MadeMap(rdbUser.Basic.Lv)
	if err != nil {
		return
	}
	rdbUser.Basic.TreasureCnt = rdbUser.Map.Treasure
	rdbUser.Basic.Dead = enums.NO
	rdbUser.Basic.DeadPointID = 0

	if err = redis.SetUser(rdbUser, locker); err != nil {
		log.Error("update redis user err", zap.Error(err), zap.Uint64("userID", rdbUser.ID))
		err = errcode.ERR_REDIS_UPDATE_USER
		return
	}

	//fmt.Printf(rdbUser.Map.String())

	msgID = pb.MsgID_ResetMap_AckId
	msg = &pb.ResetMapAck{
		Strength:      int32(rdbUser.Basic.Strength),
		Maps:          redis.BuildPbMap(rdbUser.Map),
		FreeResetMap:  rdbUser.Basic.FreeResetMap == enums.YES,
		ExtraStrength: int32(rdbUser.Basic.ExtraStrength),
	}

	return
}
