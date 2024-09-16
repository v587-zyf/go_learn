package global

import (
	"github.com/v587-zyf/gc/iface"
)

type Config struct {
	SID int64 `yaml:"sid"`
	Dev bool  `yaml:"dev"`

	Http HttpConfig `yaml:"http"`
	Tcp  TcpConfig  `yaml:"tcp"`

	//Redis rdb.ConfigSingle `yaml:"redis"`
	//Mongo db.Config        `yaml:"mongo"`

	TableDb string `yaml:"tabledb"`
}

type HttpConfig struct {
	ListenAddr string `yaml:"listen_addr" mapstructure:"listen_addr"` // 监听地址
	LinkAddr   string `yaml:"link_addr" mapstructure:"link_addr"`     // 连接地址
	Pem        string `yaml:"pem" mapstructure:"pem"`
	Key        string `yaml:"key" mapstructure:"key"`
}

type TcpConfig struct {
	ListenAddr string                  `yaml:"listen_addr" mapstructure:"listen_addr"` // 监听地址
	LinkAddr   string                  `yaml:"link_addr" mapstructure:"link_addr"`     // 连接地址
	Method     iface.ITpcSessionMethod `yaml:"-"`
}
