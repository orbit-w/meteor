package logger

import (
	"fmt"
	"github.com/golang/glog"
	"go.uber.org/zap"
)

type PrefixLogger struct {
	prefix string
}

func NewPrefixLogger(prefix string) *PrefixLogger {
	return &PrefixLogger{
		prefix: prefix,
	}
}

func (l *PrefixLogger) Info(args ...any) {
	args = append([]any{"[" + string(l.prefix) + "] "}, args...)
	glog.InfoDepth(d, args...)
}

func (l *PrefixLogger) Infof(format string, args ...any) {
	format = l.prefix + format
	glog.InfoDepth(d, fmt.Sprintf(format, args...))
}

func (l *PrefixLogger) Debugf(format string, args ...any) {
	format = l.prefix + format
	glog.WarningDepth(d, fmt.Sprintf(format, args...))
}

func (l *PrefixLogger) Error(args ...any) {
	args = append([]any{"[" + string(l.prefix) + "] "}, args...)
	glog.ErrorDepth(d, args...)
}

func (l *PrefixLogger) Errorf(format string, args ...any) {
	format = l.prefix + format
	glog.ErrorDepth(d, fmt.Sprintf(format, args...))
}

func (l *PrefixLogger) InfoDepth(depth int, args ...any) {
	glog.InfoDepth(d+depth, args...)
}

type ZapLogger struct {
	prefix string
}

func NewLogger(prefix string) *ZapLogger {
	return &ZapLogger{
		prefix: prefix,
	}
}

func (l *ZapLogger) Info(msg string, field ...zap.Field) {
	msg = "[" + string(l.prefix) + "] " + msg
	getBaseLogger().Info(msg, field...)
}

func (l *ZapLogger) Infof(format string, args ...any) {
	format = "[" + string(l.prefix) + "] " + format
	getBaseLogger().Info(fmt.Sprintf(format, args...))
}

func (l *ZapLogger) Error(msg string, field ...zap.Field) {
	msg = "[" + string(l.prefix) + "] " + msg
	getBaseLogger().Error(msg, field...)
}

func (l *ZapLogger) Errorf(format string, args ...any) {
	format = "[" + string(l.prefix) + "] " + format
	getBaseLogger().Error(fmt.Sprintf(format, args...))
}
