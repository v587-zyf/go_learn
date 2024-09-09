package session

type SessionMgr interface {
	Init(addr string, options ...string) error
	CreateSession() (Session, error)
	Get(sessionId string) (Session, error)
}
