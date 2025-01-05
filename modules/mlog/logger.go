package mlog

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/natefinch/lumberjack"
)

type Logger struct {
	logger *zap.Logger
	config *Config
	atom   zap.AtomicLevel
	sugar  *zap.SugaredLogger
}

// New creates a new logger with the given options
func New(opts ...Option) (*Logger, error) {
	l := &Logger{
		config: DefaultConfig(),
		atom:   zap.NewAtomicLevel(),
	}

	// Apply options
	for _, opt := range opts {
		opt(l)
	}

	// Validate config
	if err := l.config.Validate(); err != nil {
		return nil, err
	}

	// Build zap logger
	zapConfig := l.config.BuildZapConfig()
	zapLogger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	l.logger = zapLogger
	l.sugar = zapLogger.Sugar()

	return l, nil
}

func NewDevelopmentLogger() *Logger {
	l := &Logger{
		atom: zap.NewAtomicLevel(),
	}

	zapConfig := newDevelopmentConfig()
	zapLogger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	l.logger = zapLogger
	l.sugar = zapLogger.Sugar()

	return l
}

// With creates a child logger with the given fields
func (l *Logger) With(fields ...zapcore.Field) *Logger {
	child := *l
	child.logger = l.logger.With(fields...)
	child.sugar = child.logger.Sugar()
	return &child
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	return l.With(zap.String("Prefix", prefix))
}

// WithContext creates a child logger with context fields
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Add trace ID if available
	if ctx != nil {
		if traceID := ctx.Value("trace_id"); traceID != nil {
			return l.With(zap.Any("trace_id", traceID))
		}
	}
	return l
}

func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugar.Fatalf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugar.Panicf(format, args...)
}

// Sugar returns a sugared logger
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// NewFileLogger creates a new logger with file output
func NewFileLogger(opts ...Option) *Logger {
	l := &Logger{
		config: DefaultFileLogConfig(),
		atom:   zap.NewAtomicLevel(),
	}

	// Apply options
	for _, opt := range opts {
		opt(l)
	}

	// Validate config
	if err := l.config.Validate(); err != nil {
		panic(err)
	}

	// 配置日志轮转Hook
	logRotate := &lumberjack.Logger{
		Filename:   l.config.OutputPaths[0],          // 日志文件路径
		MaxSize:    l.config.RotateConfig.MaxSize,    // 单个日志文件最大尺寸，单位MB
		MaxBackups: l.config.RotateConfig.MaxBackups, // 保留的旧日志文件最大数量
		MaxAge:     l.config.RotateConfig.MaxAge,     // 保留的旧日志文件最大天数
		Compress:   l.config.RotateConfig.Compress,   // 是否压缩旧日志文件
	}

	// 设置日志级别
	if level, err := zapcore.ParseLevel(l.config.Level); err == nil {
		l.atom.SetLevel(level)
	} else {
		l.atom.SetLevel(zapcore.InfoLevel)
	}

	// 创建logger选项
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	// 如果是开发模式，添加开发模式选项
	if l.config.Development {
		options = append(options, zap.Development())
	}

	// 添加初始字段
	if len(l.config.InitialFields) > 0 {
		fields := make([]zap.Field, 0, len(l.config.InitialFields))
		for k, v := range l.config.InitialFields {
			fields = append(fields, zap.Any(k, v))
		}
		options = append(options, zap.Fields(fields...))
	}

	// 创建logger
	l.logger = zap.New(l.newTee(logRotate), options...)
	l.sugar = l.logger.Sugar()

	return l
}

func (l *Logger) newTee(logRotate *lumberjack.Logger) zapcore.Core {
	// 创建编码器配置
	encoderConfig := newEncoderConfig()

	// 创建控制台编码器（带颜色）和文件编码器（JSON格式）
	var consoleEncoder, fileEncoder zapcore.Encoder
	if l.config.Format == "json" {
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		// 文件始终使用JSON格式
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建Core
	var cores []zapcore.Core

	// 如果是开发模式，添加控制台输出
	if l.config.Development {
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			l.atom,
		))
	}

	// 添加文件输出
	cores = append(cores, zapcore.NewCore(
		fileEncoder,
		zapcore.AddSync(logRotate),
		l.atom,
	))
	return zapcore.NewTee(cores...)
}
