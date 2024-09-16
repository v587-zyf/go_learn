package global

import (
	"github.com/v587-zyf/gc/gcnet/ws_server"
	"github.com/v587-zyf/gc/utils"
	"os"
)

var (
	Conf      *Config
	WsServer  *ws_server.WsServer
	Snowflake *utils.Snowflake

	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
