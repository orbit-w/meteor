package logger

import (
	"fmt"
	"github.com/golang/glog"
)

type PrefixLogger struct {
	prefix string
}

func NewLogger(prefix string) *PrefixLogger {
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
