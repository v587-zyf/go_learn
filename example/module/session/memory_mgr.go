package session

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type MemorySessionMgr struct {
	sessionMap map[string]Session
	rwlock     sync.RWMutex
}

func NewMemorySessionMgr() *MemorySessionMgr {
	m := &MemorySessionMgr{
		sessionMap: make(map[string]Session, 1024),
	}

	return m
}

func (m *MemorySessionMgr) Init(addr string, options ...string) error {

	return nil
}
func (m *MemorySessionMgr) CreateSession() (Session, error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	u := uuid.New()
	sessionId := u.String()
	session := NewMemorySession(sessionId)

	return session, nil
}
func (m *MemorySessionMgr) Get(sessionId string) (Session, error) {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	session, ok := m.sessionMap[sessionId]
	if !ok {
		return nil, errors.New("session not found")
	}

	return session, nil
}
