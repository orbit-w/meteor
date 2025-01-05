package mlog

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "default config",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "with level",
			opts: []Option{
				WithLevel("debug"),
			},
			wantErr: false,
		},
		{
			name: "with invalid level",
			opts: []Option{
				WithLevel("invalid"),
			},
			wantErr: false, // Should not error as invalid level is ignored
		},
		{
			name: "with format",
			opts: []Option{
				WithFormat("json"),
			},
			wantErr: false,
		},
		{
			name: "with development mode",
			opts: []Option{
				WithDevelopment(true),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.opts...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, logger)
			assert.NotNil(t, logger.logger)
			assert.NotNil(t, logger.sugar)
		})
	}
}

func Test_NewPrint(t *testing.T) {
	logger, err := New()
	require.NoError(t, err)
	assert.NotNil(t, logger)
	logger.Error("test error")
	logger.Sync()
}

func TestLogger_WithContext(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	logger := &Logger{
		logger: zap.New(core),
	}

	// Test with trace ID
	ctx := context.WithValue(context.Background(), "trace_id", "test-trace-id")
	loggerWithCtx := logger.WithContext(ctx)
	loggerWithCtx.Info("test message")

	// Parse the log output
	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	// Verify trace ID is present
	assert.Equal(t, "test-trace-id", logEntry["trace_id"])
}

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	logger := &Logger{
		logger: zap.New(core),
		sugar:  zap.New(core).Sugar(),
	}

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		message  string
		level    string
		wantLogs bool
	}{
		{
			name:     "debug level",
			logFunc:  logger.Debugf,
			message:  "debug message",
			level:    "debug",
			wantLogs: true,
		},
		{
			name:     "info level",
			logFunc:  logger.Infof,
			message:  "info message",
			level:    "info",
			wantLogs: true,
		},
		{
			name:     "warn level",
			logFunc:  logger.Warnf,
			message:  "warn message",
			level:    "warn",
			wantLogs: true,
		},
		{
			name:     "error level",
			logFunc:  logger.Errorf,
			message:  "error message",
			level:    "error",
			wantLogs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message)

			if tt.wantLogs {
				assert.True(t, buf.Len() > 0)
				var logEntry map[string]interface{}
				err := json.Unmarshal(buf.Bytes(), &logEntry)
				require.NoError(t, err)
				assert.Equal(t, tt.level, logEntry["level"])
				assert.Equal(t, tt.message, logEntry["msg"])
			} else {
				assert.Equal(t, 0, buf.Len())
			}
		})
	}
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	logger := &Logger{
		logger: zap.New(core),
		sugar:  zap.New(core).Sugar(),
	}

	// Test with additional fields
	childLogger := logger.With(zap.String("key", "value"))
	childLogger.Info("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)
	assert.Equal(t, "value", logEntry["key"])
}

func TestLogger_WithPrefix(t *testing.T) {
	logger := NewDevelopmentLogger()

	childLogger := logger.WithPrefix("OrbitNetwork")
	assert.NotNil(t, childLogger)
	childLogger.Info("test message", zap.String("Id", "1"))
	logger.Sync()
}

func TestWithSampling(t *testing.T) {
	logger, err := New(WithSampling(1, 100, time.Second))
	require.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.config.Sampling)
	assert.Equal(t, 1, logger.config.Sampling.Initial)
	assert.Equal(t, 100, logger.config.Sampling.Thereafter)
	assert.Equal(t, time.Second, logger.config.Sampling.Tick)
}

func TestWithInitialFields(t *testing.T) {
	fields := map[string]interface{}{
		"app":     "test",
		"version": "1.0.0",
	}
	logger, err := New(WithInitialFields(fields))
	require.NoError(t, err)
	assert.NotNil(t, logger)
	assert.Equal(t, fields, logger.config.InitialFields)
}

func TestNewDevelopmentLogger(t *testing.T) {
	logger := NewDevelopmentLogger()
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.logger)
	assert.NotNil(t, logger.sugar)
	assert.NotNil(t, logger.atom)
}

func TestStop(t *testing.T) {
	assert.NotPanics(t, func() {
		Stop()
	})
}

func TestWithPrefix(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	logger := &Logger{
		logger: zap.New(core),
		sugar:  zap.New(core).Sugar(),
	}

	prefixedLogger := logger.With(zap.String("Prefix", "TestPrefix"))
	prefixedLogger.Info("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)
	assert.Equal(t, "TestPrefix", logEntry["Prefix"])
}

func TestWithPrefix2(t *testing.T) {
	prefixedLogger := WithPrefix("TestPrefix")
	prefixedLogger.Info("test message")
	prefixedLogger.Sync()
}

func TestNewFileLogger(t *testing.T) {
	logger := NewFileLogger(
		WithOutputPaths("logs/test.log"),
		WithRotation(100, 7, 10, false),
		WithLevel("error"),
	)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.sugar)
	assert.NotNil(t, logger.atom)

	logger.Info("test message")
	logger.Error("test error")

	logger.Sync()
}

func TestFileLoggerDevelopment(t *testing.T) {
	logger := NewFileLogger(
		WithDevelopment(true),
		WithOutputPaths("logs/test.log"),
		WithRotation(100, 7, 10, false),
		WithLevel("error"),
	)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.sugar)
	assert.NotNil(t, logger.atom)
	testCaller(t, logger)
	logger.Sync()
}

func TestNewFileLoggerWithConsole(t *testing.T) {
	logger := NewFileLogger(
		WithDevelopment(true),
		WithOutputPaths("logs/test.log"),
		WithRotation(100, 7, 10, false),
		WithLevel("error"),
		WithFormat("console"),
	)
	assert.NotNil(t, logger)
	assert.NotNil(t, logger.sugar)
	assert.NotNil(t, logger.atom)

	logger.Info("test message")
	logger.Error("test error")
	logger.Sync()
}

func testCaller(t *testing.T, logger *Logger) {
	logger.Errorf("test message")
	testCaller2(t, logger)
}

func testCaller2(_ *testing.T, logger *Logger) {
	logger.Errorf("test message")
}
