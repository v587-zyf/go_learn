package registry

import "time"

type Options struct {
	// 地址
	Addrs []string
	// 超时
	Timeout time.Duration
	// 心跳时间
	Heartbeat int64
	// 注册地址
	RegistryPath string
}

// 函数类型变量
type Option func(opts *Options)

func WithAddrs(addrs []string) Option {
	return func(opts *Options) {
		opts.Addrs = addrs
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *Options) {
		opts.Timeout = timeout
	}
}

func WithHeartbeat(heartbeat int64) Option {
	return func(opts *Options) {
		opts.Heartbeat = heartbeat
	}
}

func WithRegistryPath(path string) Option {
	return func(opts *Options) {
		opts.RegistryPath = path
	}
}
