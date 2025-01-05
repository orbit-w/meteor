package mlog

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Option configures a Logger.
type Option func(*Logger)

// WithLevel sets the logger's level.
func WithLevel(level string) Option {
	return func(l *Logger) {
		if lvl, err := zapcore.ParseLevel(level); err == nil {
			l.config.Level = level
			l.atom.SetLevel(lvl)
		}
	}
}

// WithFormat sets the logger's encoding format.
func WithFormat(format string) Option {
	return func(l *Logger) {
		if format == "json" || format == "console" {
			l.config.Format = format
		}
	}
}

// WithOutputPaths sets the logger's output paths.
func WithOutputPaths(paths ...string) Option {
	return func(l *Logger) {
		l.config.OutputPaths = paths
	}
}

// WithDevelopment sets the logger's development mode.
func WithDevelopment(development bool) Option {
	return func(l *Logger) {
		l.config.Development = development
	}
}

// WithSampling sets the logger's sampling config.
func WithSampling(initial, thereafter int, tick time.Duration) Option {
	return func(l *Logger) {
		l.config.Sampling = &SamplingConfig{
			Initial:    initial,
			Thereafter: thereafter,
			Tick:       tick,
		}
	}
}

// WithInitialFields sets initial fields for the logger.
func WithInitialFields(fields map[string]interface{}) Option {
	return func(l *Logger) {
		l.config.InitialFields = fields
	}
}

// WithRotation sets the logger's rotation config.
func WithRotation(maxSize, maxAge, maxBackups int, compress bool) Option {
	return func(l *Logger) {
		l.config.RotateConfig = &RotateConfig{
			MaxSize:    maxSize,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
			Compress:   compress,
		}
	}
}
