package global

type Config struct {
	SID int64 `yaml:"sid"`

	Grpc GrpcConfig `yaml:"grpc"`

	Redis RedisConfig `yaml:"redis"`
	Mongo MongoConfig `yaml:"mongo"`

	Table Table `yaml:"table"`

	Dev   bool `yaml:"dev"`
	Debug bool `yaml:"debug"`

	Telegram TelegramConfig `yaml:"telegram"`
}

type GrpcConfig struct {
	ListenAddr string `yaml:"listen_addr" mapstructure:"listen_addr"` // 监听地址
	LinkAddr   string `yaml:"link_addr" mapstructure:"link_addr"`     // 连接地址
}

type RedisConfig struct {
	Addrs []string `yaml:"addrs" mapstructure:"addrs"`
	Addr  string   `yaml:"addr" mapstructure:"addr"`
	Pass  string   `yaml:"pass" mapstructure:"pass"`
}

type MongoConfig struct {
	Uri string `yaml:"uri" mapstructure:"uri"`
	DB  string `yaml:"db" mapstructure:"db"`
	Use bool   `yaml:"use" mapstructure:"use"`
}

type Table struct {
	Path string `yaml:"path" mapstructure:"path"`
}

type TelegramUnit struct {
	Name    string  `yaml:"name" mapstructure:"name"`
	Token   string  `yaml:"token" mapstructure:"token"`
	ChatIDs []int64 `yaml:"chat_ids" mapstructure:"chat_ids"`
}
type TelegramConfig struct {
	LoginToken string          `yaml:"login_token" mapstructure:"login_token"`
	ListenDev  []*TelegramUnit `yaml:"listen_dev" mapstructure:"listen_dev"`
	Listen     []*TelegramUnit `yaml:"listen" mapstructure:"listen"`
}
