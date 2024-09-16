package db

import (
	pb "comm/t_proto/out/client"
	"fmt"
	"github.com/v587-zyf/gc/iface"
	"sync"
)

type Point struct {
	Id     int  `json:"id,omitempty" bson:"id,omitempty"`
	X      int  `json:"x,omitempty" bson:"x,omitempty"`
	Y      int  `json:"y,omitempty" bson:"y,omitempty"`
	IsSafe bool `json:"isSafe,omitempty" bson:"isSafe,omitempty"`

	Thunder  int `json:"thunder,omitempty" bson:"thunder,omitempty"`   // 周围地雷数量
	Treasure int `json:"treasure,omitempty" bson:"treasure,omitempty"` // 宝箱金币

	Objs  map[iface.TileType]struct{} `json:"objs,omitempty" bson:"objs,omitempty"`
	oLock sync.RWMutex
}

func NewGrid(id int, x, y int) *Point {
	p := &Point{
		Id:   id,
		X:    x,
		Y:    y,
		Objs: make(map[iface.TileType]struct{}),
	}

	return p
}

func (p *Point) IsSafed() bool  { return p.IsSafe }
func (p *Point) SetSafe(i bool) { p.IsSafe = i }
func (p *Point) ID() int        { return p.Id }
func (p *Point) AddObj(t iface.TileType) *Point {
	p.oLock.Lock()
	defer p.oLock.Unlock()

	p.Objs[t] = struct{}{}

	return p
}
func (p *Point) RemoveObj(t iface.TileType) *Point {
	p.oLock.Lock()
	defer p.oLock.Unlock()

	delete(p.Objs, t)

	return p
}
func (p *Point) GetObjs() map[iface.TileType]struct{} {
	p.oLock.RLock()
	defer p.oLock.RUnlock()

	return p.Objs
}
func (p *Point) GetObjsSlice() []iface.TileType {
	p.oLock.RLock()
	defer p.oLock.RUnlock()

	slice := make([]iface.TileType, len(p.Objs))

	i := 0
	for tileType := range p.Objs {
		slice[i] = tileType
	}

	return slice
}
func (p *Point) HasObj(t iface.TileType) bool {
	p.oLock.RLock()
	defer p.oLock.RUnlock()

	if _, ok := p.Objs[t]; ok {
		return true
	}

	return false
}

func (p *Point) SetThunder(t int)  { p.Thunder = t }
func (p *Point) GetThunder() int   { return p.Thunder }
func (p *Point) SetTreasure(t int) { p.Treasure = t }
func (p *Point) GetTreasure() int  { return p.Treasure }

func (p *Point) IsWallOrEmpty() bool {
	p.oLock.RLock()
	defer p.oLock.RUnlock()

	//if p.HasObj(iface.Thunder) || p.HasObj(iface.Player) || p.HasObj(iface.Treasure) {
	if p.HasObj(iface.Thunder) || p.HasObj(iface.Treasure) {
		return false
	}

	return true
}

func (p *Point) ToPb() *pb.MapUnit {
	pbUnit := &pb.MapUnit{
		Id:      int32(p.ID()),
		X:       int32(p.X),
		Y:       int32(p.Y),
		Thunder: int32(p.Thunder),
	}
	switch {
	case p.HasObj(iface.Player):
		pbUnit.TileType = pb.TileType_Player
	case p.HasObj(iface.Wall):
		pbUnit.TileType = pb.TileType_Wall
	case p.HasObj(iface.Treasure):
		pbUnit.TileType = pb.TileType_Treasure
	case p.HasObj(iface.Thunder):
		pbUnit.TileType = pb.TileType_Thunder
	default:
		pbUnit.TileType = pb.TileType_Empty
	}

	return pbUnit
}

func (p *Point) String() string {
	return fmt.Sprintf("Id:%d X:%d Y:%d IsSafed:%t Objs:%v Thunder:%d Treasure:%d\n",
		p.Id, p.X, p.Y, p.IsSafed(), p.Objs, p.Thunder, p.Treasure)
}
