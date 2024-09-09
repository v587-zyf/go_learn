package session

var (
	sessionMgr SessionMgr
)

func Init(provider string, addr string, options ...string) error {
	switch provider {
	case "memory":
		sessionMgr = NewMemorySessionMgr()
	case "redis":
		sessionMgr = NewRedisSessionMgr()
	}

	if err := sessionMgr.Init(addr, options...); err != nil {
		return err
	}

	return nil
}
