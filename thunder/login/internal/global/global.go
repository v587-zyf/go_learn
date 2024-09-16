package global

import (
	"github.com/v587-zyf/gc/gcnet/http_server"
	"github.com/v587-zyf/gc/utils"
	"os"
)

var (
	Conf      *Config
	Snowflake *utils.Snowflake

	HttpServer *http_server.HttpServer

	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
