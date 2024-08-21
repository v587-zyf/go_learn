package log

import (
	"go.uber.org/zap/zapcore"
)

const (
	DefaultMaxSize    = 300 * 1024 * 1024
	DefaultMaxAge     = 30
	DefaultMaxBackups = 999

	DefaultInfoFileName  = "%s/info-%d-%d-%d.log"
	DefaultErrorFileName = "%s/error-%d-%d-%d.log"
)

type Option struct {
	serName string
	serID   int64

	level        zapcore.Level
	isStdout     bool // output console
	isStackTrace bool // output stack
	skipCaller   int

	infoPath    string
	infoMaxSize int // file max size(byte)
	infoMaxAge  int // file timestamps(days)

	errPath    string
	errMaxSize int // file max size(MB) default=100MB
	errMaxAge  int // file timestamps(day)
}

type OptionFn func(opt *Option)

func NewOption() *Option {
	return &Option{
		serName: "",
		serID:   0,

		level:        zapcore.DebugLevel,
		isStdout:     true,
		isStackTrace: true,
		skipCaller:   3,

		infoPath:    "./log",
		infoMaxSize: DefaultMaxSize,
		infoMaxAge:  DefaultMaxAge,

		errPath:    "./log",
		errMaxSize: DefaultMaxSize,
		errMaxAge:  DefaultMaxAge,
	}
}

func WithSerName(data string) OptionFn {
	return func(opt *Option) {
		opt.serName = data
	}
}

func WithSerID(data int64) OptionFn {
	return func(opt *Option) {
		opt.serID = data
	}
}

func WithLevel(data zapcore.Level) OptionFn {
	return func(opt *Option) {
		opt.level = data
	}
}

func WithIsStdout(data bool) OptionFn {
	return func(opt *Option) {
		opt.isStdout = data
	}
}

func WithIsStackTrace(data bool) OptionFn {
	return func(opt *Option) {
		opt.isStackTrace = data
	}
}

func WithSkipCaller(data int) OptionFn {
	return func(opt *Option) {
		opt.skipCaller = data
	}
}

func WithInfoPath(data string) OptionFn {
	return func(opt *Option) {
		opt.infoPath = data
	}
}

func WithInfoMaxSize(data int) OptionFn {
	return func(opt *Option) {
		opt.infoMaxSize = data
	}
}

func WithInfoMaxAge(data int) OptionFn {
	return func(opt *Option) {
		opt.infoMaxAge = data
	}
}

func WithErrPath(data string) OptionFn {
	return func(opt *Option) {
		opt.errPath = data
	}
}

func WithErrMaxSize(data int) OptionFn {
	return func(opt *Option) {
		opt.errMaxSize = data
	}
}

func WithErrMaxAge(data int) OptionFn {
	return func(opt *Option) {
		opt.errMaxAge = data
	}
}
