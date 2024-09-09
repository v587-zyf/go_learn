package global

import (
	"kernel/server/http_server"
	"kernel/utils"
	"os"
)

var (
	Conf      *Config
	Snowflake *utils.Snowflake

	HttpServer *http_server.HttpServer

	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
