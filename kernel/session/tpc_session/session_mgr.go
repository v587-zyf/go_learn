package tpc_session

import (
	"kernel/iface"
	"sync"
)

var sessionMgr *SessionMgr

func init() {
	sessionMgr = NewSessionMgr()
}

type SessionMgr struct {
	online    map[uint64]iface.ITcpSession
	onlineMux sync.RWMutex
}

func GetSessionMgr() *SessionMgr {
	return sessionMgr
}

func NewSessionMgr() *SessionMgr {
	s := &SessionMgr{
		online: make(map[uint64]iface.ITcpSession),
	}

	return s
}

func (s *SessionMgr) Length() int {
	s.onlineMux.RLock()
	defer s.onlineMux.RUnlock()

	return len(s.online)
}

func (s *SessionMgr) GetOne(UID uint64) iface.ITcpSession {
	s.onlineMux.RLock()
	c, ok := s.online[UID]
	s.onlineMux.RUnlock()
	if !ok {
		return nil
	}

	return c
}

func (s *SessionMgr) IsOnline(UID uint64) bool {
	s.onlineMux.RLock()
	defer s.onlineMux.RUnlock()

	_, ok := s.online[UID]

	return ok
}

func (s *SessionMgr) Add(ss iface.ITcpSession) {
	s.onlineMux.Lock()
	defer s.onlineMux.Unlock()

	SID := ss.GetID()
	s.online[SID] = ss
}

func (s *SessionMgr) Disconnect(SID uint64) {
	s.onlineMux.Lock()
	delete(s.online, SID)
	s.onlineMux.Unlock()
}

func (s *SessionMgr) Once(UID uint64, fn func(SS iface.ITcpSession)) {
	cli := s.GetOne(UID)
	if cli == nil {
		fn(nil)
		return
	}

	fn(cli)
}

func (s *SessionMgr) Range(fn func(SID uint64, SS iface.ITcpSession)) {
	s.onlineMux.RLock()
	defer s.onlineMux.RUnlock()

	for SID, SS := range s.online {
		fn(SID, SS)
	}
}

func (s *SessionMgr) Close() {
	s.onlineMux.RLock()
	defer s.onlineMux.RUnlock()

	for _, SS := range s.online {
		SS.Close()
	}
}
