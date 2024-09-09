package session

import (
	"errors"
	"sync"
)

type MemorySession struct {
	sessionId string
	data      map[string]any
	dataMu    sync.RWMutex
}

func NewMemorySession(sessionId string) *MemorySession {
	s := &MemorySession{
		sessionId: sessionId,
		data:      make(map[string]any, 16),
	}

	return s
}

var _ = MemorySession{}

func (m *MemorySession) Set(key string, value any) error {
	m.dataMu.Lock()
	defer m.dataMu.Unlock()

	m.data[key] = value

	return nil
}
func (m *MemorySession) Get(key string) (any, error) {
	m.dataMu.RLock()
	defer m.dataMu.RUnlock()

	val, ok := m.data[key]
	if !ok {
		return nil, errors.New("session: key not found")
	}

	return val, nil
}
func (m *MemorySession) Del(key string) error {
	m.dataMu.Lock()
	defer m.dataMu.Unlock()

	delete(m.data, key)

	return nil
}
func (m *MemorySession) Save() error {

	return nil
}
