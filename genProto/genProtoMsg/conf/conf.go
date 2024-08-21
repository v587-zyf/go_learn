package conf

type Config struct {
	Source string
	Out    string
}

var config *Config

func Init(s, o string) {
	config = &Config{
		Source: s,
		Out:    o,
	}
}

func GetConf() *Config {
	return config
}
