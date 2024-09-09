package example

import (
	"go.uber.org/zap"
	"kernel/log"
)

func Log() {
	log.Info("this is Info log", zap.Int("1", 1))
	log.Error("this is Error log", zap.Int("1", 1))
	log.Debug("this is Debug log", zap.Int("1", 1))
	log.Warn("this is Warn log", zap.Int("1", 1))
}
