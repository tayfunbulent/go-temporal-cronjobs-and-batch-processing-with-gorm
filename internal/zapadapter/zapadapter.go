package zapadapter

import (
	"go.temporal.io/sdk/log"
	"go.uber.org/zap"
)

type ZapAdapter struct {
	logger *zap.Logger
}

func NewZapAdapter(logger *zap.Logger) log.Logger {
	return &ZapAdapter{logger: logger}
}

func (z *ZapAdapter) Debug(message string, keyValues ...interface{}) {
	z.logger.Sugar().Debugw(message, keyValues...)
}

func (z *ZapAdapter) Info(message string, keyValues ...interface{}) {
	z.logger.Sugar().Infow(message, keyValues...)
}

func (z *ZapAdapter) Warn(message string, keyValues ...interface{}) {
	z.logger.Sugar().Warnw(message, keyValues...)
}

func (z *ZapAdapter) Error(message string, keyValues ...interface{}) {
	z.logger.Sugar().Errorw(message, keyValues...)
}
