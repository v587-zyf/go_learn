package beego_log

import "github.com/astaxie/beego/logs"

type logger struct {
	*logs.BeeLogger
}

var defaultLog *logger
var defaultLogLevel = logs.LevelInformational

func DefaultLoggerInit() {
	defaultLog = &logger{BeeLogger: Get("default", true)}
	defaultLog.SetLogFuncCallDepth(3)
	defaultLogLevel = GetLogLevel("default")
}

func Debug(format string, v ...interface{}) {
	defaultLog.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	defaultLog.Info(format, v)
}

func Warn(format string, v ...interface{}) {
	defaultLog.Warn(format, v)
}

func Error(format string, v ...interface{}) {
	defaultLog.Error(format, v)
}
