package session

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"sync"
)

type RedisSession struct {
	sessionId string
	pool      *redis.Pool

	// 设置session 可以先放内存的map中
	// 批量导入redis 提升性能
	sessionMap map[string]any
	rwlock     sync.RWMutex

	flag int
}

const (
	SessionFlagNone = iota
	SessionFlagModify
)

func NewRedisSession(sessionId string, pool *redis.Pool) *RedisSession {
	s := &RedisSession{
		sessionId:  sessionId,
		pool:       pool,
		sessionMap: make(map[string]any, 16),
		flag:       SessionFlagNone,
	}

	return s
}

var _ = RedisSession{}

func (s *RedisSession) Set(key string, value any) error {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	s.sessionMap[key] = value
	s.flag = SessionFlagModify

	return nil
}
func (s *RedisSession) Get(key string) (any, error) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	if v, ok := s.sessionMap[key]; ok {
		return v, nil
	}

	return nil, errors.New("key not exists")
}
func (s *RedisSession) Del(key string) error {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	delete(s.sessionMap, key)
	s.flag = SessionFlagModify

	return nil
}
func (s *RedisSession) Save() error {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()

	if s.flag == SessionFlagNone {
		return nil
	}

	data, err := json.Marshal(s.sessionMap)
	if err != nil {
		return err
	}

	conn := s.pool.Get()

	_, err = conn.Do("SET", s.sessionId, string(data))
	if err != nil {
		return err
	}

	return nil
}
func (s *RedisSession) loadFromRedis() error {
	conn := s.pool.Get()
	reply, err := conn.Do("GET", s.sessionId)
	if err != nil {
		return err
	}

	if reply == nil {
		return errors.New("session not exists")
	}

	data, err := redis.String(reply, err)
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(data), &s.sessionMap); err != nil {
		return err
	}

	return nil
}
