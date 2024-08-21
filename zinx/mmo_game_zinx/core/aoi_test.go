package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	// 初始化aoi
	aoiMgr := NewAOIManager(0,
		250,
		5,
		0,
		250,
		5)

	// 打印
	fmt.Println(aoiMgr)
}

func TestGetSurroundGridByGid(t *testing.T) {
	// 初始化aoi
	aoiMgr := NewAOIManager(0,
		250,
		5,
		0,
		250,
		5)

	for gid := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridByGid(gid)
		fmt.Println("gid:", gid, " gridsLen:", len(grids))
		GIDS := make([]int, 0, len(grids))
		for _, grid := range grids {
			GIDS = append(GIDS, grid.GID)
		}
		fmt.Println("surrounding grid IDS:", GIDS)
	}
}
