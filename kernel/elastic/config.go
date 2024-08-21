package elastic

import "fmt"

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Use  bool   `json:"use"`
}

func (this *Config) ToHostString() string {
	return fmt.Sprintf("%s:%s", this.Host, this.Port)
}
