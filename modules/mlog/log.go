package mlog

import (
	"fmt"

	"go.uber.org/zap"
)

func Error(msg string, fields ...zap.Field) {
	baseLogger.Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	baseLogger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	baseLogger.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	baseLogger.Warn(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	baseLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	baseLogger.Panic(msg, fields...)
}

func Errorf(format string, args []any, field ...zap.Field) {
	msg := fmt.Sprintf(format, args...)
	baseLogger.Error(msg, field...)
}
