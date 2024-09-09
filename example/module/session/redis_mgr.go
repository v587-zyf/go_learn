package session

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"sync"
	"time"
)

type RedisSessionMgr struct {
	addr     string
	password string
	pool     *redis.Pool

	sessionMap map[string]Session
	rwlock     sync.RWMutex
}

func myPool(addr, password string) *redis.Pool {
	p := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			_, err = conn.Do("AUTH", password)
			if err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
		MaxIdle:         64,
		MaxActive:       1000,
		IdleTimeout:     240 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}

	return p
}

func NewRedisSessionMgr() *RedisSessionMgr {
	m := &RedisSessionMgr{
		sessionMap: make(map[string]Session, 32),
	}

	return m
}

func (m *RedisSessionMgr) Init(addr string, options ...string) error {
	m.addr = addr
	if len(options) > 0 {
		m.password = options[0]
	}
	m.pool = myPool(addr, m.password)

	return nil
}
func (m *RedisSessionMgr) CreateSession() (Session, error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	u := uuid.New()
	sessionId := u.String()
	session := NewRedisSession(sessionId, m.pool)

	return session, nil
}
func (m *RedisSessionMgr) Get(sessionId string) (Session, error) {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	session, ok := m.sessionMap[sessionId]
	if !ok {
		return nil, errors.New("session not found")
	}

	return session, nil
}
