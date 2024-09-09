package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogLv uint16

type Logger interface {
	Debug(str string, a ...any)
	Trace(str string, a ...any)
	Info(str string, a ...any)
	Warn(str string, a ...any)
	Error(str string, a ...any)
	Fatal(str string, a ...any)
}

const (
	DEBUG LogLv = iota
	TRACE
	INFO
	WARN
	ERROR
	FATAL
)

func str2Lv(level string) LogLv {
	lv := strings.ToLower(level)
	switch lv {
	case "debug":
		return DEBUG
	case "trace":
		return TRACE
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return DEBUG
	}
}

func lv2Str(lv LogLv) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

func logPrint(lv LogLv, format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	funcName, fileName, line := getInfo(3)
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n",
		getTime(), lv2Str(lv),
		fileName, funcName, line, str)
}

func getTime() string {
	timeN := time.Now()
	return timeN.Format("2006-01-02 15:04:05")
}

func getInfo(skip int) (funcName, fileName string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	return
}
