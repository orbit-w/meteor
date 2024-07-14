package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(fileName string, lv zapcore.Level) *zap.Logger {
	encoder := newEncoder()
	core := zapcore.NewCore(encoder, zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     14,    // days
		Compress:   false, // disabled by default
		LocalTime:  true,
	}), lv)

	return zap.New(core)
}

// Stop flushing any buffered log entries
// Applications should take care to call Sync before exiting
func Stop(log *zap.Logger) {
	if log != nil {
		_ = log.Sync()
	}
}

func newEncoder() zapcore.Encoder {
	c := zap.NewProductionEncoderConfig()
	c.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	c.EncodeName = zapcore.FullNameEncoder
	return zapcore.NewConsoleEncoder(c)
}
