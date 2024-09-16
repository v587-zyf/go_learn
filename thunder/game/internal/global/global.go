package global

import (
	"github.com/v587-zyf/gc/utils"
	"os"
)

var (
	Conf      *Config
	Snowflake *utils.Snowflake

	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
