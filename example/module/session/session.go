package session

type Session interface {
	Set(key string, value any) error
	Get(key string) (any, error)
	Del(key string) error
	Save() error
}
