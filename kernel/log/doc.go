package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defLog *Logger

func GetDefaultLogger() *Logger {
	return defLog
}

func Init(ctx context.Context, opts ...OptionFn) (err error) {
	defLog = NewLogger()
	if err = defLog.Init(ctx, opts...); err != nil {
		return err
	}
	return err
}

func Info(msg string, fields ...zapcore.Field) {
	defLog.Info(msg, fields...)
}

func Debug(msg string, fields ...zapcore.Field) {
	defLog.Debug(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	defLog.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	defLog.Error(msg, fields...)
}

func With(fields ...zap.Field) {
	defLog.With(fields...)
}
