package example

import (
	"github.com/v587-zyf/gc/log"
	"go.uber.org/zap"
)

func Log() {
	log.Info("this is Info log", zap.Int("1", 1))
	log.Error("this is Error log", zap.Int("1", 1))
	log.Debug("this is Debug log", zap.Int("1", 1))
	log.Warn("this is Warn log", zap.Int("1", 1))
}
