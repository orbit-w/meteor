package mlog

import (
	"fmt"

	"go.uber.org/zap"
)

func Error(msg string, fields ...zap.Field) {
	getBaseLogger().Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	getBaseLogger().Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	getBaseLogger().Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getBaseLogger().Warn(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	getBaseLogger().DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	getBaseLogger().Panic(msg, fields...)
}

func Errorf(format string, args []any, field ...zap.Field) {
	msg := fmt.Sprintf(format, args...)
	getBaseLogger().Error(msg, field...)
}

func Infof(format string, args []any, field ...zap.Field) {
	msg := fmt.Sprintf(format, args...)
	getBaseLogger().Info(msg, field...)
}

func Debugf(format string, args []any, field ...zap.Field) {
	msg := fmt.Sprintf(format, args...)
	getBaseLogger().Debug(msg, field...)
}

func Warnf(format string, args []any, field ...zap.Field) {
	msg := fmt.Sprintf(format, args...)
	getBaseLogger().Warn(msg, field...)
}

func DPanicf(format string, args []any, field ...zap.Field) {
	msg := fmt.Sprintf(format, args...)
	getBaseLogger().DPanic(msg, field...)
}

func Panicf(format string, args []any, field ...zap.Field) {
	msg := fmt.Sprintf(format, args...)
	getBaseLogger().Panic(msg, field...)
}
