package global

type Config struct {
	SID int64 `yaml:"sid"`

	Dev bool `yaml:"dev"`

	Http HttpConfig `yaml:"http"`
	Grpc GrpcConfig `yaml:"grpc"`

	Redis RedisConfig `yaml:"redis"`
	Mongo MongoConfig `yaml:"mongo"`

	Table Table `yaml:"table"`

	Telegram TelegramConfig `yaml:"telegram"`
}

type HttpConfig struct {
	Https        bool   `yaml:"https" mapstructure:"https"`
	ListenAddr   string `yaml:"listen_addr" mapstructure:"listen_addr"`
	LinkAddr     string `yaml:"link_addr" mapstructure:"link_addr"`
	Key          string `yaml:"key" mapstructure:"key"`
	Pem          string `yaml:"pem" mapstructure:"pem"`
	AllowOrigins string `yaml:"allowOrigins" mapstructure:"allowOrigins"`
}

type GrpcConfig struct {
	ListenAddr string `yaml:"listen_addr" mapstructure:"listen_addr"`
	LinkAddr   string `yaml:"link_addr" mapstructure:"link_addr"`
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

type TelegramConfig struct {
	LoginToken   string `yaml:"login_token" mapstructure:"login_token"`
	ClientUrl    string `yaml:"client_url" mapstructure:"client_url"`
	StartPhoto   string `yaml:"start_photo" mapstructure:"start_photo"`
	StartCaption string `yaml:"start_caption" mapstructure:"start_caption"`
}
