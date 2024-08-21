package core

import "fmt"

/*
 * X轴坐标：idx Y轴坐标：idy
 * 格子编号： id = idy * cntX + idx(利用格子坐标得到格子编号)
 * 格子坐标： idx = id % cntX ... idy = id / cntX(利用格子编号得到格子坐标)
 * 格子X轴坐标： idx = id % nx(利用格子id得到X轴坐标编号)
 * 格子Y轴坐标： idy = id / nx(利用格子id得到Y轴坐标编号)
 */

const (
	AOI_MIN_X = 85
	AOI_MAX_X = 410
	AOI_CNT_X = 10

	AOI_MIN_Y = 75
	AOI_MAX_Y = 400
	AOI_CNT_Y = 20
)

type AOIManager struct {
	MinX  int           // 左边界
	MaxX  int           // 右边界
	CntX  int           // X轴格子数量
	MinY  int           // 上边界
	MaxY  int           // 下边界
	CntY  int           // Y轴格子数量
	grids map[int]*Grid // 格子ID：格子对象
}

func NewAOIManager(minX, maxX, cntX, minY, maxY, cntY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntX:  cntX,
		MinY:  minY,
		MaxY:  maxY,
		CntY:  cntY,
		grids: make(map[int]*Grid),
	}

	// 给所有格子进行编号和初始化
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			// 格子编号: id = y * cntX + x
			gid := y*cntX + x

			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}

	return aoiMgr
}

// 得到每个格子在X轴方向的宽度
func (aoiMgr *AOIManager) gridWidth() int {
	return (aoiMgr.MaxX - aoiMgr.MinX) / aoiMgr.CntX
}

// 得到每个格子在Y轴方向的长度
func (aoiMgr *AOIManager) gridLength() int {
	return (aoiMgr.MaxY - aoiMgr.MinY) / aoiMgr.CntY
}

// 试调专用-打印格子信息
func (aoiMgr *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX=%d MaxX=%d CntX=%d MinY=%d MaxY=%d CntY=%d\n",
		aoiMgr.MinX, aoiMgr.MaxX, aoiMgr.CntX, aoiMgr.MinY, aoiMgr.MaxY, aoiMgr.CntY)

	for _, grid := range aoiMgr.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// 根据各自GID得到周边格子ID集合 九宫格
func (aoiMgr *AOIManager) GetSurroundGridByGid(gID int) (grids []*Grid) {
	// 当前id是否在AOIManager中
	nowGrid, ok := aoiMgr.grids[gID]
	if !ok {
		return
	}

	// 初始化grids返回切片 将当前格子放入九宫格中
	grids = append(grids, nowGrid)

	// 判断GID左边是否有格子 右边是否有格子
	// 通过GID得到当前格子X轴编号 idx = id % nx
	idx := gID % aoiMgr.CntX

	// 判断idx左边是否还有格子 有就放在gidsX中
	if idx > 0 {
		grids = append(grids, aoiMgr.grids[gID-1])
	}

	// 判断idx右边是否还有格子 有就放在gIdsX中
	if idx < aoiMgr.CntX-1 {
		grids = append(grids, aoiMgr.grids[gID+1])
	}

	// 把当前九宫格中X轴格子放一个集合中
	gIdsX := make([]int, 0, len(grids))
	for _, grid := range grids {
		gIdsX = append(gIdsX, grid.GID)
	}

	// 遍历gIdsX每个格子的gid
	for _, gid := range gIdsX {
		// 得到当前格子Y轴编号 idy = id / ny
		idy := gid / aoiMgr.CntY
		// gid上边是否有格子
		if idy > 0 {
			grids = append(grids, aoiMgr.grids[gid-aoiMgr.CntX])
		}
		// gid下边是否有格子
		if idy < aoiMgr.CntY-1 {
			grids = append(grids, aoiMgr.grids[gid+aoiMgr.CntX])
		}
	}

	return
}

func (aoiMgr *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - aoiMgr.MinX) / aoiMgr.gridWidth()
	idy := (int(y) - aoiMgr.MinY) / aoiMgr.gridLength()

	return idy*aoiMgr.CntX + idx
}

// 根据X坐标获取九宫格
func (aoiMgr *AOIManager) GetElemsSByPos(x, y float32) (elems []int) {
	// 得到当前格子id
	gID := aoiMgr.GetGidByPos(x, y)
	// 通过GID得到周边九宫格
	grids := aoiMgr.GetSurroundGridByGid(gID)
	// 将九宫格信息累加到elems
	for _, grid := range grids {
		elems = append(elems, grid.GetElements()...)
		//fmt.Println("--->gID:", grid.GID, " pIDs:", grid.GetElements())
	}

	return
}

// 添加一个element到格子中
func (aoiMgr *AOIManager) AddEleToGrid(eID, gID int) {
	aoiMgr.grids[gID].AddElement(eID)
}

// 删除格子中一个element
func (aoiMgr *AOIManager) RemoveFromEidAndGid(eID, gID int) {
	aoiMgr.grids[gID].RemoveElement(eID)
}

// 通过GID获取全部ElementIDs
func (aoiMgr *AOIManager) GetEleIdsByGid(gID int) (eIDs []int) {
	eIDs = aoiMgr.grids[gID].GetElements()

	return
}

// 通过坐标将Element添加到格子
func (aoiMgr *AOIManager) AddToGridByPos(eID int, x, y float32) {
	gID := aoiMgr.GetGidByPos(x, y)
	grid := aoiMgr.grids[gID]
	grid.AddElement(eID)
}

// 通过坐标删除格子中一个Element
func (aoiMgr *AOIManager) RemoveFromPos(eID int, x, y float32) {
	gID := aoiMgr.GetGidByPos(x, y)
	grid := aoiMgr.grids[gID]
	grid.RemoveElement(eID)
}
