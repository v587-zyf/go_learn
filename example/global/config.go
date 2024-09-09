package global

type Config struct {
	SID int64 `yaml:"sid"`

	Http HttpConfig `yaml:"http"`

	Elastic ElasticConfig `yaml:"elastic"`
	Redis   RedisConfig   `yaml:"redis"`
	Mongo   MongoConfig   `yaml:"mongo"`

	Telegram TelegramConfig `yaml:"telegram"`
}

type HttpConfig struct {
	ListenAddr string `yaml:"listen_addr" mapstructure:"listen_addr"` // 监听地址
	LinkAddr   string `yaml:"link_addr" mapstructure:"link_addr"`     // 连接地址
}

type RedisConfig struct {
	Addrs []string `yaml:"addrs" mapstructure:"addrs"`
	Addr  string   `yaml:"addr" mapstructure:"addr"`
	Pass  string   `yaml:"pass" mapstructure:"pass"`
}

type MongoConfig struct {
	Use bool   `yaml:"use" mapstructure:"use"`
	Uri string `yaml:"uri" mapstructure:"uri"`
	DB  string `yaml:"db" mapstructure:"db"`
}

type ElasticConfig struct {
	Use  bool   `yaml:"use" mapstructure:"use"`
	Host string `yaml:"host" mapstructure:"host"`
	Port string `yaml:"port" mapstructure:"port"`
}

type TelegramConfig struct {
	Use   bool   `yaml:"use" mapstructure:"use"`
	Token string `yaml:"token" mapstructure:"token"`
}
