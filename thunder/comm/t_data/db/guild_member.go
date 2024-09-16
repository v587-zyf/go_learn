package db

import (
	"sync"
)

type GuildMembersUnit struct{}
type GuildMembers struct {
	Data map[uint64]*GuildMembersUnit `bson:"data" json:"data"`
}

var guildMemLock sync.RWMutex

func NewGuildMembers() *GuildMembers {
	return &GuildMembers{
		Data: make(map[uint64]*GuildMembersUnit),
	}
}

func (g *GuildMembers) Has(UID uint64) bool {
	guildMemLock.RLock()
	defer guildMemLock.RUnlock()

	_, ok := g.Data[UID]
	if !ok {
		return false
	}

	return true
}

func (g *GuildMembers) Add(UID uint64) {
	guildMemLock.Lock()
	defer guildMemLock.Unlock()

	g.Data[UID] = &GuildMembersUnit{}
}

func (g *GuildMembers) Del(UID uint64) {
	guildMemLock.Lock()
	defer guildMemLock.Unlock()

	delete(g.Data, UID)
}

func (g *GuildMembers) Len() int {
	guildMemLock.RLock()
	defer guildMemLock.RUnlock()

	return len(g.Data)
}
