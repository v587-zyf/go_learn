package global

type Config struct {
	SID int64 `yaml:"sid"`

	Grpc GrpcConfig `yaml:"grpc"`
}

type GrpcConfig struct {
	ListenAddr string `yaml:"listen_addr" mapstructure:"listen_addr"` // 监听地址
	LinkAddr   string `yaml:"link_addr" mapstructure:"link_addr"`     // 连接地址
}
