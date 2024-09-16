package db

import (
	enum "comm/t_enum"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
	"math/rand"
)

/*
 * X轴坐标：idx Y轴坐标：idy
 * 格子编号： Id = idy * cntX + idx(利用格子坐标得到格子编号)
 * 格子坐标： idx = Id % cntX ... idy = Id / cntX(利用格子编号得到格子坐标)
 * 格子X轴坐标： idx = Id % nx(利用格子id得到X轴坐标编号)
 * 格子Y轴坐标： idy = Id / nx(利用格子id得到Y轴坐标编号)
 */

type Map struct {
	ID int `bson:"id,omitempty" json:"id,omitempty"`

	MinX int `bson:"minX,omitempty" json:"minX,omitempty"` // 左边界
	MaxX int `bson:"maxX,omitempty" json:"maxX,omitempty"` // 右边界
	CntX int `bson:"cntX,omitempty" json:"cntX,omitempty"` // X轴格子数量

	MinY int `bson:"minY,omitempty" json:"minY,omitempty"` // 上边界
	MaxY int `bson:"maxY,omitempty" json:"maxY,omitempty"` // 下边界
	CntY int `bson:"cntY,omitempty" json:"cntY,omitempty"` // Y轴格子数量

	CntId    int `bson:"cntId,omitempty" json:"cntId,omitempty"`     // 总id
	BirthId  int `bson:"birthId,omitempty" json:"birthId,omitempty"` // 出生点
	NowId    int `bson:"nowId" json:"nowId,omitempty"`               // 当前玩家所在点
	Treasure int `bson:"treasure" json:"treasure,omitempty"`         // 宝箱数量

	Grids map[int]*Point `bson:"grids,omitempty" json:"grids,omitempty"` // 格子ID：格子对象
}

func NewMap(minX, maxX, cntX, minY, maxY, cntY int) *Map {
	m := &Map{
		MinX: minX,
		MaxX: maxX,
		CntX: cntX,

		MinY: minY,
		MaxY: maxY,
		CntY: cntY,

		Grids: make(map[int]*Point),
	}

	var gid int
	// 给所有格子进行编号和初始化
	for y := 0; y < cntY; y++ {
		for x := 0; x < cntX; x++ {
			// 格子编号: Id = Y * cntX + X
			gid = y*cntX + x

			m.Grids[gid] = NewGrid(gid,
				m.MinX+x*m.gridWidth(),
				//m.MinX+(x+1)*m.gridWidth(),
				m.MinY+y*m.gridLength())
			//m.MinY+(y+1)*m.gridLength())
			m.Grids[gid].AddObj(iface.Wall)
		}
	}
	m.CntId = gid

	return m
}

// 得到每个格子在X轴方向的宽度
func (m *Map) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntX
}

// 得到每个格子在Y轴方向的长度
func (m *Map) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntY
}

// 根据GID得到周边格子ID集合 九宫格
func (m *Map) GetSurroundGridByGid(gID int) (grids []*Point) {
	// 当前id是否在AOIManager中
	nowGrid, ok := m.Grids[gID]
	if !ok {
		log.Error("gID err", zap.Int("gid", gID))
		return
	}

	// 初始化grids返回切片 将当前格子放入九宫格中
	grids = append(grids, nowGrid)

	// 判断GID左边是否有格子 右边是否有格子
	// 通过GID得到当前格子X轴编号 idx = Id % nx
	idx := gID % m.CntX

	// 判断idx左边是否还有格子 有就放在gidsX中
	if idx > 0 {
		grids = append(grids, m.Grids[gID-1])
	}

	// 判断idx右边是否还有格子 有就放在gIdsX中
	if idx < m.CntX-1 {
		grids = append(grids, m.Grids[gID+1])
	}

	// 把当前九宫格中X轴格子放一个集合中
	gIdsX := make([]int, 0, len(grids))
	for _, grid := range grids {
		gIdsX = append(gIdsX, grid.ID())
	}

	// 遍历gIdsX每个格子的gid
	for _, gid := range gIdsX {
		// 得到当前格子Y轴编号 idy = Id / ny
		idy := gid / m.CntY
		// gid上边是否有格子
		if idy > 0 {
			grids = append(grids, m.Grids[gid-m.CntX])
		}
		// gid下边是否有格子
		if idy < m.CntY-1 {
			grids = append(grids, m.Grids[gid+m.CntX])
		}
	}

	return
}

// 根据各自GID得到周边格子ID集合 九宫格
func (m *Map) GetSurroundGridsByGid(gID int) (grids map[enum.MAP_TOWARD]*Point) {
	// 当前id是否在AOIManager中
	nowGrid, ok := m.Grids[gID]
	if !ok {
		log.Error("gID err", zap.Int("gid", gID))
		return
	}
	grids[enum.MAP_NOW] = nowGrid

	idx := gID % m.CntX
	if idx > 0 {
		grids[enum.MAP_LEFT] = m.Grids[gID-1]
	}
	if idx < m.CntX-1 {
		grids[enum.MAP_RIGHT] = m.Grids[gID+1]
	}

	var (
		gid     int
		_toward enum.MAP_TOWARD
	)
	for toward, point := range grids {
		gid = point.ID()
		// 得到当前格子Y轴编号 idy = Id / ny
		idy := gid / m.CntY
		// gid上边是否有格子
		if idy > 0 {
			switch toward {
			case enum.MAP_NOW:
				_toward = enum.MAP_TOP
			case enum.MAP_LEFT:
				_toward = enum.MAP_LEFT_TOP
			case enum.MAP_RIGHT:
				_toward = enum.MAP_RIGHT_TOP
			default:
				panic("unhandled default case")
			}
			grids[_toward] = m.Grids[gid-m.CntX]
		}
		// gid下边是否有格子
		if idy < m.CntY-1 {
			switch toward {
			case enum.MAP_NOW:
				_toward = enum.MAP_BOTTOM
			case enum.MAP_LEFT:
				_toward = enum.MAP_RIGHT_BOTTOM
			case enum.MAP_RIGHT:
				_toward = enum.MAP_LEFT_BOTTOM
			default:
				panic("unhandled default case")
			}
			grids[_toward] = m.Grids[gid+m.CntX]
		}
	}

	return
}

// 根据各自GID得到周边格子ID集合 上下左右
func (m *Map) GetAllSidesGridByGid(gID int) (grids []*Point) {
	// 当前id是否在AOIManager中
	nowGrid, ok := m.Grids[gID]
	if !ok {
		log.Error("gID err", zap.Int("gid", gID))
		return
	}
	// self
	grids = append(grids, nowGrid)

	idx := m.GetXByGid(gID)
	// left
	if idx > 0 {
		grids = append(grids, m.Grids[gID-1])
	}
	// right
	if idx < m.CntX-1 {
		grids = append(grids, m.Grids[gID+1])
	}

	idy := m.GetYByGid(gID)
	// top
	if idy > 0 {
		grids = append(grids, m.Grids[gID-m.CntX])
	}
	// bottom
	if idy < m.CntY-1 {
		grids = append(grids, m.Grids[gID+m.CntX])
	}

	return
}

func (m *Map) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntX + idx
}

func (m *Map) GetGridByPos(x, y int32) *Point {
	return m.Grids[m.GetGidByPos(float32(x), float32(y))]
}

// 根据X坐标获取九宫格
func (m *Map) GetObjsSByPos(x, y float32) map[iface.TileType]int {
	// 得到当前格子id
	gID := m.GetGidByPos(x, y)
	// 通过GID得到周边九宫格
	grids := m.GetSurroundGridByGid(gID)

	tileMap := make(map[iface.TileType]int)
	// 将九宫格信息累加到objs
	for _, grid := range grids {
		for tileType := range grid.GetObjs() {
			tileMap[tileType]++
		}
	}

	return tileMap
}

// 添加一个element到格子中
func (m *Map) AddObjToGrid(gid int, t iface.TileType) {
	m.Grids[gid].AddObj(t)
}

// 删除格子中一个obj
func (m *Map) RemoveFromEidAndGid(gid int, t iface.TileType) {
	m.Grids[gid].RemoveObj(t)
}

// 通过GID获取全部obj
func (m *Map) GetEleIdsByGid(gid int) map[iface.TileType]struct{} {
	return m.Grids[gid].GetObjs()
}

func (m *Map) GetGridByGid(gid int) *Point {
	return m.Grids[gid]
}

// 通过坐标将obj添加到格子
func (m *Map) AddToGridByPos(t iface.TileType, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.Grids[gID]
	grid.AddObj(t)
}

// 通过坐标删除格子中一个obj
func (m *Map) RemoveFromPos(t iface.TileType, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.Grids[gID]
	grid.RemoveObj(t)
}

func (m *Map) GetXByGid(gid int) int {
	return gid % m.CntX
}

func (m *Map) GetYByGid(gid int) int {
	return gid / m.CntX
}

func (m *Map) SetBirthGid(gid int) { m.BirthId = gid }
func (m *Map) SetBirthXY(x, y float32) {
	gid := m.GetGidByPos(x, y)
	m.BirthId = gid
}
func (m *Map) GetBirthGid() int { return m.BirthId }
func (m *Map) GetBirthXY() (x, y float32) {
	gid := m.GetBirthGid()
	x = float32(m.GetXByGid(gid))
	y = float32(m.GetYByGid(gid))

	return
}

func (m *Map) SetNowGid(gid int) { m.NowId = gid }
func (m *Map) SetNowXY(x, y float32) {
	gid := m.GetGidByPos(x, y)
	m.NowId = gid
}
func (m *Map) GetNowGid() int { return m.NowId }
func (m *Map) GetNowXY() (x, y float32) {
	gid := m.GetNowGid()
	x = float32(m.GetXByGid(gid))
	y = float32(m.GetYByGid(gid))
	return
}

func (m *Map) GetTreasure() int { return m.Treasure }
func (m *Map) SetTreasure(t int) {
	m.Treasure = t
	return
}

func (m *Map) RandBirthPos() int {
	return rand.Intn(m.CntId)
}
func (m *Map) RandThunderPos() (gid int) {
	var (
		grid *Point
	)

LOOP:
	for {
		gid = rand.Intn(m.CntId)
		grid = m.GetGridByGid(gid)
		// 宝箱与地雷不重叠 && 不能在安全区
		if !grid.IsWallOrEmpty() || m.GetGridByGid(gid).IsSafed() {
			continue
		}
		break LOOP
	}

	return
}
func (m *Map) RandTreasurePos() (gid int) {
	var (
		grid     *Point
		allSides []*Point
	)

LOOP:
	for {
		gid = rand.Intn(m.CntId)
		grid = m.GetGridByGid(gid)
		// 宝箱与地雷不重叠
		if grid.HasObj(iface.Thunder) || grid.HasObj(iface.Player) || grid.HasObj(iface.Treasure) {
			continue
		}
		// 宝箱周围是4格需至少有一个是空地
		allSides = m.GetAllSidesGridByGid(gid)
		flag := false
		for _, side := range allSides {
			if !side.IsWallOrEmpty() {
				continue
			}
			flag = true
		}
		if flag {
			break LOOP
		}
	}
	return
}

func (m *Map) IsSurroundHasThunder(gID int) (result bool, grids []*Point) {
	grids = m.GetSurroundGridByGid(gID)
	for _, grid := range grids {
		if grid.HasObj(iface.Thunder) {
			return true, nil
		}
	}

	return
}
func (m *Map) GetSurroundNoThunderPids(point *Point, filterMap map[int]struct{}) map[int]struct{} {
	if _, ok := filterMap[point.ID()]; ok || point.HasObj(iface.Thunder) || point.HasObj(iface.Empty) || point.HasObj(iface.Player) {
		return filterMap
	}

	filterMap[point.ID()] = struct{}{}
	ret, grids := m.IsSurroundHasThunder(point.ID())
	if !ret {
		for _, pit := range grids {
			for k := range m.GetSurroundNoThunderPids(pit, filterMap) {
				filterMap[k] = struct{}{}
			}
		}
	}

	return filterMap
}

// 根据各自GID得到周边格子ID集合 二十五个格子
func (m *Map) GetSurroundGrids25(gID int) (grids []*Point) {
	// 当前id是否在AOIManager中
	nowGrid, ok := m.Grids[gID]
	if !ok {
		log.Error("gID err", zap.Int("gid", gID))
		return
	}

	grids = append(grids, nowGrid)

	idx := gID % m.CntX
	// left
	if idx > 0 {
		grids = append(grids, m.Grids[gID-1])
		if idx-1 > 0 {
			grids = append(grids, m.Grids[gID-2])
		}
	}
	// right
	if idx < m.CntX-1 {
		grids = append(grids, m.Grids[gID+1])
		if idx+1 < m.CntX-1 {
			grids = append(grids, m.Grids[gID+2])
		}
	}

	gIdsX := make([]int, 0, len(grids))
	for _, grid := range grids {
		gIdsX = append(gIdsX, grid.ID())
	}

	for _, gid := range gIdsX {
		idy := gid / m.CntY
		// top
		if idy > 0 {
			grids = append(grids, m.Grids[gid-m.CntX])
			if idy-1 > 0 {
				grids = append(grids, m.Grids[gid-m.CntX*2])
			}
		}
		// bottom
		if idy < m.CntY-1 {
			grids = append(grids, m.Grids[gid+m.CntX])
			if idy+1 < m.CntY-1 {
				grids = append(grids, m.Grids[gid+m.CntX*2])
			}
		}
	}

	return
}

func (m *Map) IsOver() bool {
	for _, point := range m.Grids {
		if point.HasObj(iface.Treasure) || (point.HasObj(iface.Wall) && point.HasObj(iface.Empty)) {
			return false
		}
	}

	return true
}

// 试调专用-打印格子信息
func (m *Map) String() string {
	//s := fmt.Sprintf("Map:\n MinX=%d MaxX=%d CntX=%d MinY=%d MaxY=%d CntY=%d\n",
	//	m.MinX, m.MaxX, m.CntX, m.MinY, m.MaxY, m.CntY)
	//
	//for _, grid := range m.Grids {
	//	s += fmt.Sprintln(grid)
	//}

	var (
		s    string
		gid  int
		grid *Point
	)
	for y := 0; y < m.CntY; y++ {
		for x := 0; x < m.CntX; x++ {
			// 格子编号: Id = Y * cntX + X
			gid = y*m.CntX + x

			grid = m.GetGridByGid(gid)
			switch {
			case grid.HasObj(iface.Treasure):
				s += fmt.Sprintf("B(id:%d.x:%d.y:%d.th:%d.tr:%d)\t", grid.ID(), grid.X, grid.Y, grid.Thunder, grid.Treasure)
			case grid.HasObj(iface.Player):
				s += fmt.Sprintf("P(id:%d.x:%d.y:%d.th:%d.tr:%d)\t", grid.ID(), grid.X, grid.Y, grid.Thunder, grid.Treasure)
			case grid.HasObj(iface.Thunder):
				s += fmt.Sprintf("T(id:%d.x:%d.y:%d.th:%d.tr:%d)\t", grid.ID(), grid.X, grid.Y, grid.Thunder, grid.Treasure)
			case grid.HasObj(iface.Wall):
				s += fmt.Sprintf("W(id:%d.x:%d.y:%d.th:%d.tr:%d)\t", grid.ID(), grid.X, grid.Y, grid.Thunder, grid.Treasure)
			case grid.HasObj(iface.Empty):
				s += fmt.Sprintf("E(id:%d.x:%d.y:%d.th:%d.tr:%d)\t", grid.ID(), grid.X, grid.Y, grid.Thunder, grid.Treasure)
			default:
				s += fmt.Sprintf("N(id:%d.x:%d.y:%d.th:%d.tr:%d)\t", grid.ID(), grid.X, grid.Y, grid.Thunder, grid.Treasure)
			}
			//if grid.HasObj(iface.Treasure) {
			//	s += fmt.Sprintf("b\t")
			//} else if grid.HasObj(iface.Player) {
			//	s += fmt.Sprintf("p\t")
			//} else if grid.HasObj(iface.Thunder) {
			//	s += fmt.Sprintf("t\t")
			//} else if grid.HasObj(iface.Wall) {
			//	s += fmt.Sprintf("w\t")
			//}
		}
		s += "\n"
	}
	s += fmt.Sprintf("birthId=%d\n", m.BirthId)
	s += fmt.Sprintf("nowId=%d\n", m.NowId)
	s += fmt.Sprintf("Treasure=%d\n", m.Treasure)

	//for _, grid = range m.Grids {
	//	s += fmt.Sprintln(grid)
	//}

	return s
}
