package rdb_single

type RedisSingleOption struct {
	addr string
	pwd  string
}

type Option func(o *RedisSingleOption)

func NewRedisSingleOption() *RedisSingleOption {
	return &RedisSingleOption{}
}

func WithAddr(addr string) Option {
	return func(o *RedisSingleOption) {
		o.addr = addr
	}
}

func WithPwd(pwd string) Option {
	return func(o *RedisSingleOption) {
		o.pwd = pwd
	}
}
