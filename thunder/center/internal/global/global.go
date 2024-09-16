package global

import (
	"github.com/robfig/cron/v3"
	"github.com/v587-zyf/gc/utils"
	"os"
)

var (
	Conf      *Config
	Snowflake *utils.Snowflake
	Cron      *cron.Cron

	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
