package mlog

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config defines the configuration for logger
type Config struct {
	// Level is the minimum enabled logging level
	Level string `yaml:"level" json:"level"`
	// Format specifies the format of the logger output
	// Valid values are "json" and "console"
	Format string `yaml:"format" json:"format"`
	// OutputPaths is a list of URLs or file paths to write logging output to
	OutputPaths []string `yaml:"outputPaths" json:"outputPaths"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally
	Development bool `yaml:"development" json:"development"`
	// Sampling sets a sampling strategy for the logger
	Sampling *SamplingConfig `yaml:"sampling" json:"sampling"`
	// InitialFields sets the initial fields for the logger
	InitialFields map[string]interface{} `yaml:"initialFields" json:"initialFields"`
	// RotateConfig sets the configuration for log rotation
	RotateConfig *RotateConfig `yaml:"rotateConfig" json:"rotateConfig"`
}

// SamplingConfig sets a sampling strategy for the logger
type SamplingConfig struct {
	Initial    int           `yaml:"initial" json:"initial"`
	Thereafter int           `yaml:"thereafter" json:"thereafter"`
	Tick       time.Duration `yaml:"tick" json:"tick"`
}

// RotateConfig sets the configuration for log rotation
type RotateConfig struct {
	MaxSize    int  `yaml:"maxSize" json:"maxSize"`       // megabytes
	MaxAge     int  `yaml:"maxAge" json:"maxAge"`         // days
	MaxBackups int  `yaml:"maxBackups" json:"maxBackups"` // files
	Compress   bool `yaml:"compress" json:"compress"`
}

// DefaultConfig returns a Config with default settings
func DefaultConfig() *Config {
	return &Config{
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout"},
		Development: false,
		Sampling: &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
			Tick:       time.Second,
		},
		InitialFields: make(map[string]interface{}),
	}
}

func DefaultFileLogConfig() *Config {
	return &Config{
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"logs/log.log"},
		Development: false,
		RotateConfig: &RotateConfig{ // 日志轮转配置
			MaxSize:    500,   // 单个文件最大 500MB
			MaxAge:     7,     // 保留 7 天
			MaxBackups: 10,    // 保留 10 个备份
			Compress:   false, // 压缩旧文件
		},
		InitialFields: make(map[string]interface{}),
	}
}

// Validate validates the config and sets default values
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Set defaults if not set
	if c.Level == "" {
		c.Level = "info"
	}
	if c.Format == "" {
		c.Format = "json"
	}
	if len(c.OutputPaths) == 0 {
		c.OutputPaths = []string{"stdout"}
	}
	if c.Sampling == nil {
		c.Sampling = &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
			Tick:       time.Second,
		}
	}
	if c.InitialFields == nil {
		c.InitialFields = make(map[string]interface{})
	}

	// Validate level
	_, err := zapcore.ParseLevel(c.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	// Validate format
	if c.Format != "json" && c.Format != "console" {
		return fmt.Errorf("invalid format: %s, must be either json or console", c.Format)
	}

	if c.RotateConfig != nil {
		if c.RotateConfig.MaxSize <= 0 {
			return fmt.Errorf("invalid max size: %d, must be greater than 0", c.RotateConfig.MaxSize)
		}
		if c.RotateConfig.MaxAge <= 0 {
			return fmt.Errorf("invalid max age: %d, must be greater than 0", c.RotateConfig.MaxAge)
		}
		if c.RotateConfig.MaxBackups < 0 {
			return fmt.Errorf("invalid max backups: %d, must be greater than or equal to 0", c.RotateConfig.MaxBackups)
		}
	} else {
		c.RotateConfig = &RotateConfig{
			MaxSize:    500,
			MaxAge:     7,
			MaxBackups: 10,
			Compress:   false,
		}
	}

	return nil
}

// BuildZapConfig converts Config to zap.Config
func (c *Config) BuildZapConfig() zap.Config {
	level, _ := zapcore.ParseLevel(c.Level)

	return zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: c.Development,
		Sampling: &zap.SamplingConfig{
			Initial:    c.Sampling.Initial,
			Thereafter: c.Sampling.Thereafter,
		},
		Encoding:         c.Format,
		EncoderConfig:    newEncoderConfig(),
		OutputPaths:      c.OutputPaths,
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    c.InitialFields,
	}
}

// newDevelopmentConfig 实例化开发环境日志
// StackTraces are included on logs of WarnLevel and above.
// Warn 级别以上，包含堆栈信息
func newDevelopmentConfig() zap.Config {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return cfg
}

// newEncoderConfig returns a new EncoderConfig for structured logging
func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
