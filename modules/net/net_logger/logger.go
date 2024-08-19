package net_logger

import (
	"github.com/orbit-w/meteor/bases/zap_logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	prefix string
	logger *zap.Logger
}

func NewProductionLogger(fileName, prefix string, lv zapcore.Level) *Logger {
	log := new(Logger)
	log.prefix = prefix
	log.logger = zap_logger.NewProductionLogger(fileName, lv)
	log.logger.With()
	return log
}

func NewDevelopmentLogger(prefix string) *Logger {
	log := new(Logger)
	log.prefix = prefix
	log.logger = zap_logger.NewDevelopmentLogger()
	return log
}

func NewLogger(prefix string) *Logger {
	log := new(Logger)
	log.prefix = prefix
	log.logger = getBaseLogger().With()
	return log
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	msg = l.prefix + msg
	l.logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	msg = l.prefix + msg
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	msg = l.prefix + msg
	l.logger.Error(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	msg = l.prefix + msg
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Stop() {
	zap_logger.Stop(l.logger)
}
