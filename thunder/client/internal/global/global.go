package global

import (
	"client/internal/robot_mgr"
	"os"
)

var (
	Conf *Config

	RobotManager *robot_mgr.RobotManager

	exitChan   = make(chan struct{})
	signalChan = make(chan os.Signal, 1)
)
