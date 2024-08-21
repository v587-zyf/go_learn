package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type Logger struct {
	options *Option

	fields     []zap.Field
	infoLogger *zap.Logger
	errLogger  *zap.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func NewLogger() *Logger {
	return &Logger{
		options: NewOption(),
		fields:  make([]zap.Field, 0),
	}
}

func (l *Logger) Init(ctx context.Context, opts ...OptionFn) (err error) {
	l.ctx, l.cancel = context.WithCancel(ctx)
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(l.options)
		}
	}
	if l.options.serName != "" {
		l.fields = append(l.fields, zap.String("serName", l.options.serName))
	}
	if l.options.serID != 0 {
		l.fields = append(l.fields, zap.Int64("serID", l.options.serID))
	}

	l.InitInfo()
	return nil
}

func (l *Logger) InitInfo() {
	now := time.Now()
	fileEncoder := zapcore.NewJSONEncoder(DefaultFileEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(DefaultConsoleEncoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), l.options.level)
	infoHook := &lumberjack.Logger{
		Filename:   fmt.Sprintf(DefaultInfoFileName, l.options.infoPath, now.Year(), now.Month(), now.Day()), // filePath
		MaxSize:    l.options.infoMaxSize,                                                                    // 单个日志文件最大大小（以MB为单位）
		MaxBackups: DefaultMaxBackups,                                                                        // 保留旧文件的最大个数
		MaxAge:     l.options.infoMaxAge,                                                                     // 保留旧文件的最大天数
		Compress:   false,                                                                                    // 是否压缩
	}
	defer infoHook.Close()

	infoFileWriteSyncer := zapcore.AddSync(infoHook)
	infoFileCore := zapcore.NewCore(fileEncoder, infoFileWriteSyncer, l.options.level)
	logCore := zapcore.NewTee(infoFileCore, consoleCore)
	l.infoLogger = zap.New(logCore, zap.AddCaller(), zap.AddCallerSkip(l.options.skipCaller), zap.Fields(l.fields...))

	errHook := &lumberjack.Logger{
		Filename:   fmt.Sprintf(DefaultErrorFileName, l.options.infoPath, now.Year(), now.Month(), now.Day()), // filePath
		MaxSize:    l.options.infoMaxSize,                                                                     // 单个日志文件最大大小（以MB为单位）
		MaxBackups: DefaultMaxBackups,                                                                         // 保留旧文件的最大个数
		MaxAge:     l.options.infoMaxAge,                                                                      // 保留旧文件的最大天数
		Compress:   false,                                                                                     // 是否压缩
	}
	defer errHook.Close()
	errFileWriteSyncer := zapcore.AddSync(errHook)
	errFileCore := zapcore.NewCore(fileEncoder, errFileWriteSyncer, l.options.level)
	logCore = zapcore.NewTee(errFileCore, infoFileCore, consoleCore)
	l.errLogger = zap.New(logCore, zap.AddCaller(), zap.AddCallerSkip(l.options.skipCaller), zap.Fields(l.fields...))

}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.infoLogger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.infoLogger.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.infoLogger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.errLogger.Error(msg, fields...)
}

func (l *Logger) With(fields ...zap.Field) {
	l.infoLogger.With(fields...)
	l.errLogger.With(fields...)
}
