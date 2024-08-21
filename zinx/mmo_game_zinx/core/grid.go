package core

import (
	"fmt"
	"sync"
)

/*
*
AOI地图中格子类型
*/
type Grid struct {
	GID     int          // 格子ID
	MinX    int          // 左边界
	MaxX    int          // 右边界
	MinY    int          // 上边界
	MaxY    int          // 下边界
	Element map[int]bool // 格子中的物体
	eleLock sync.RWMutex // 格子物体锁
}

func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	g := &Grid{
		GID:     gID,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    minY,
		MaxY:    maxY,
		Element: make(map[int]bool),
	}
	return g
}

// 给格子添加物体
func (g *Grid) AddElement(id int) {
	g.eleLock.Lock()
	defer g.eleLock.Unlock()

	g.Element[id] = true
}

// 给格子删除物体
func (g *Grid) RemoveElement(id int) {
	g.eleLock.Lock()
	defer g.eleLock.Unlock()

	delete(g.Element, id)
}

// 获取格子所有元素
func (g *Grid) GetElements() []int {
	g.eleLock.RLock()
	defer g.eleLock.RUnlock()

	elements := make([]int, 0)
	for key := range g.Element {
		elements = append(elements, key)
	}

	return elements
}

// 试调专用-打印格子元素
func (g *Grid) String() string {
	return fmt.Sprintf("GID:%d MinX:%d MaxX:%d MinY:%d MaxY:%d Element:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.Element)
}
