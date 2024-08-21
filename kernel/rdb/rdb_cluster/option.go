package rdb_cluster

type RedisClusterOption struct {
	addrs []string
	pwd   string
}

type Option func(o *RedisClusterOption)

func NewRedisClusterOption() *RedisClusterOption {
	return &RedisClusterOption{}
}

func WithAddr(addrs []string) Option {
	return func(o *RedisClusterOption) {
		o.addrs = addrs
	}
}

func WithPwd(pwd string) Option {
	return func(o *RedisClusterOption) {
		o.pwd = pwd
	}
}
