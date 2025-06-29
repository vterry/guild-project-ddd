package logger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	l, _ := zap.NewProduction()
	return &ZapLogger{logger: l.Sugar()}
}

func (z *ZapLogger) Info(msg string, args ...interface{})  { z.logger.Infow(msg, args...) }
func (z *ZapLogger) Warn(msg string, args ...interface{})  { z.logger.Warnw(msg, args...) }
func (z *ZapLogger) Error(msg string, args ...interface{}) { z.logger.Errorw(msg, args...) }
func (z *ZapLogger) Debug(msg string, args ...interface{}) { z.logger.Debugw(msg, args...) }
