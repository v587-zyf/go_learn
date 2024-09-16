package global

import (
	"github.com/v587-zyf/gc/gcnet/grpc_server"
	"os"
)

var (
	Conf *Config

	GrpcServer *grpc_server.GrpcServer

	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
