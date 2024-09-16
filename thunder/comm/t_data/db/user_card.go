package db

import (
	"comm/t_tdb"
	"sync"
)

type Card struct {
	Data  map[int]int `bson:"data" json:"data,omitempty"`
	NoNum int         `bson:"no_num" json:"no_num,omitempty"` // 未获得次数 每次获得卡牌概率30% 连续3次未获得，第四次必得
}

var cardLock sync.RWMutex

func NewCard() *Card {
	return &Card{
		Data: map[int]int{
			t_tdb.Conf().Init_card_id: 1,
		}}
}

func (c *Card) Has(id int) bool {
	cardLock.RLock()
	defer cardLock.RUnlock()

	if _, ok := c.Data[id]; ok {
		return true
	}

	return false
}

func (c *Card) Reset() {
	cardLock.Lock()
	defer cardLock.Unlock()

	c.Data[t_tdb.Conf().Init_card_id] = 1
}

func (c *Card) AddLv(id int) {
	cardLock.Lock()
	defer cardLock.Unlock()

	c.Data[id]++
}

func (c *Card) SetLv(id, lv int) {
	cardLock.Lock()
	defer cardLock.Unlock()

	c.Data[id] = lv
}

func (c *Card) GetLv(id int) int {
	cardLock.Lock()
	defer cardLock.Unlock()

	return c.Data[id]
}

func (c *Card) Add(id int) {
	cardLock.Lock()
	defer cardLock.Unlock()

	c.Data[id] = 1
}

func (c *Card) Remove(id int) {
	cardLock.Lock()
	defer cardLock.Unlock()

	delete(c.Data, id)
}

func (c *Card) GetNum() int {
	return c.NoNum
}
func (c *Card) AddNum() {
	c.NoNum++
}
func (c *Card) ResetNum() {
	c.NoNum = 0
}

func (c *Card) ToPb() map[int32]int32 {
	cardLock.RLock()
	defer cardLock.RUnlock()

	ids := make(map[int32]int32)
	for id, lv := range c.Data {
		ids[int32(id)] = int32(lv)
	}

	return ids
}
