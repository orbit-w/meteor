package logger

import (
	"github.com/orbit-w/meteor/bases/zap_logger"
	"go.uber.org/zap"
)

var (
	baseLogger *zap.Logger
)

const (
	FlagLogToStderr = "alsologtostderr" //是否输出到stderr
	FlagLogDir      = "log_dir"         //日志存储路径
	FlagV           = "v"               //日志等级
)

const d = 2

func getBaseLogger() *zap.Logger {
	if baseLogger == nil {
		baseLogger = zap_logger.NewDevelopmentLogger()
	}
	return baseLogger
}

func SetBaseLogger(logger *zap.Logger) {
	if logger == nil {
		panic("global logger invalid")
	}
	baseLogger = logger
}

func StopBaseLogger() {
	if baseLogger != nil {
		_ = baseLogger.Sync()
	}
}
