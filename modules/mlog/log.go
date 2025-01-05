package mlog

import (
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var (
	// global is the global logger instance
	global = NewDevelopmentLogger()
)

// Global logger methods
func Info(msg string, fields ...zap.Field)      { global.Info(msg, fields...) }
func Infof(format string, args ...interface{})  { global.Infof(format, args...) }
func Debug(msg string, fields ...zap.Field)     { global.Debug(msg, fields...) }
func Debugf(format string, args ...interface{}) { global.Debugf(format, args...) }
func Warn(msg string, fields ...zap.Field)      { global.Warn(msg, fields...) }
func Warnf(format string, args ...interface{})  { global.Warnf(format, args...) }
func Error(msg string, fields ...zap.Field)     { global.Error(msg, fields...) }
func Errorf(format string, args ...interface{}) { global.Errorf(format, args...) }
func Fatal(msg string, fields ...zap.Field)     { global.Fatal(msg, fields...) }
func Fatalf(format string, args ...interface{}) { global.Fatalf(format, args...) }
func Panic(msg string, fields ...zap.Field)     { global.Panic(msg, fields...) }
func Panicf(format string, args ...interface{}) { global.Panicf(format, args...) }

func With(fields ...zap.Field) *Logger {
	return global.With(fields...)
}

func WithPrefix(prefix string) *Logger {
	return global.With(zap.String("Prefix", prefix))
}

func Stop() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Second))
	defer func() {
		cancel()
	}()

	ch := make(chan struct{}, 1)

	go func() {
		_ = global.Sync()
		close(ch)
	}()

	select {
	case <-ch:
	case <-ctx.Done():

	}
}
