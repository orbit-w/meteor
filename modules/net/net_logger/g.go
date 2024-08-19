package net_logger

import (
	"github.com/orbit-w/meteor/bases/zap_logger"
	"go.uber.org/zap"
)

var (
	baseLogger *zap.Logger
)

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
