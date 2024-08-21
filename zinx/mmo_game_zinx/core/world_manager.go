package core

import (
	"sync"
)

type WorldManager struct {
	// AOIManager 世界地图AOI管理模块
	AOIManager *AOIManager
	// 当前在线Players集合
	Players map[int32]*Player
	// Players保护锁
	pLock sync.RWMutex
}

var WorldMgr *WorldManager

func init() {
	WorldMgr = &WorldManager{
		AOIManager: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_Y),
		Players:    make(map[int32]*Player),
	}
}

// 添加1个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	defer wm.pLock.Unlock()
	wm.Players[player.Pid] = player

	wm.AOIManager.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

// 删除1个玩家
func (wm *WorldManager) RemovePlayer(pId int32) {
	wm.pLock.Lock()
	defer wm.pLock.Unlock()

	player, ok := wm.Players[pId]
	if !ok {
		return
	}
	wm.AOIManager.RemoveFromPos(int(pId), player.X, player.Z)

	delete(wm.Players, pId)
}

// 通过玩家ID查询Player
func (wm *WorldManager) GetPlayerByPid(pId int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pId]
}

// 获取全部在线玩家
func (wm *WorldManager) GetAllPlayer() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)
	for _, player := range wm.Players {
		players = append(players, player)
	}

	return players
}
