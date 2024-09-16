package global

import "github.com/v587-zyf/gc/iface"

type Config struct {
	SID int64 `yaml:"sid"`
	Dev bool  `yaml:"dev"`

	Grpc GrpcConfig `yaml:"grpc"`
	Ws   WsConfig   `yaml:"ws"`
	Tcp  TcpConfig  `yaml:"tcp"`

	Redis RedisConfig `yaml:"redis"`
}

type WsConfig struct {
	Https    bool   `yaml:"https" mapstructure:"https"`
	Addr     string `yaml:"addr" mapstructure:"addr"`
	LinkAddr string `yaml:"addr" mapstructure:"link_addr"`
	Pem      string `yaml:"pem" mapstructure:"pem"`
	Key      string `yaml:"key" mapstructure:"key"`
}

type GrpcConfig struct {
	ListenAddr string `yaml:"listen_addr" mapstructure:"listen_addr"` // 监听地址
	LinkAddr   string `yaml:"link_addr" mapstructure:"link_addr"`     // 连接地址
}

type TcpConfig struct {
	ListenAddr string                  `yaml:"listen_addr" mapstructure:"listen_addr"` // 监听地址
	LinkAddr   string                  `yaml:"link_addr" mapstructure:"link_addr"`     // 连接地址
	Method     iface.ITpcSessionMethod `yaml:"-"`
}

type RedisConfig struct {
	Addrs []string `yaml:"addrs" mapstructure:"addrs"`
	Addr  string   `yaml:"addr" mapstructure:"addr"`
	Pass  string   `yaml:"pass" mapstructure:"pass"`
}

type Table struct {
	Path string `yaml:"path" mapstructure:"path"`
}
