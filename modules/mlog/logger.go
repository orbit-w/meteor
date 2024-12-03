package mlog

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type ZapLogger struct {
	prefix string
}

func NewLogger(prefix string) *ZapLogger {
	return &ZapLogger{
		prefix: prefix,
	}
}

func (l *ZapLogger) Info(msg string, field ...zap.Field) {
	msg = strings.Join([]string{"[", l.prefix, "] ", msg}, "")
	getBaseLogger().Info(msg, field...)
}

func (l *ZapLogger) Infof(format string, args []any, field ...zap.Field) {
	format = strings.Join([]string{"[", l.prefix, "] ", format}, "")
	getBaseLogger().Info(fmt.Sprintf(format, args...), field...)
}

func (l *ZapLogger) Debug(msg string, field ...zap.Field) {
	msg = strings.Join([]string{"[", l.prefix, "] ", msg}, "")
	getBaseLogger().Debug(msg, field...)
}

func (l *ZapLogger) Debugf(format string, args []any, field ...zap.Field) {
	format = strings.Join([]string{"[", l.prefix, "] ", format}, "")
	getBaseLogger().Debug(fmt.Sprintf(format, args...), field...)
}

func (l *ZapLogger) Error(msg string, field ...zap.Field) {
	msg = strings.Join([]string{"[", l.prefix, "] ", msg}, "")
	getBaseLogger().Error(msg, field...)
}

func (l *ZapLogger) Errorf(format string, args []any, field ...zap.Field) {
	format = strings.Join([]string{"[", l.prefix, "] ", format}, "")
	getBaseLogger().Error(fmt.Sprintf(format, args...), field...)
}

func (l *ZapLogger) Warn(msg string, field ...zap.Field) {
	msg = strings.Join([]string{"[", l.prefix, "] ", msg}, "")
	getBaseLogger().Warn(msg, field...)
}

func (l *ZapLogger) Warnf(format string, args []any, field ...zap.Field) {
	format = strings.Join([]string{"[", l.prefix, "] ", format}, "")
	getBaseLogger().Warn(fmt.Sprintf(format, args...), field...)
}
