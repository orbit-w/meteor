package prefix_logger

import (
	"github.com/golang/glog"
	"github.com/orbit-w/meteor/bases/zap_logger"
	"go.uber.org/zap"
)

type Logger struct {
	prefix string
}

func (l *Logger) Info(args ...any) {
	msg = l.prefix + msg
	glog.Info()
	l.logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	msg = l.prefix + msg
	glog.
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
