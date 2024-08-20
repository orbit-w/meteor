package logger

import (
	"github.com/orbit-w/meteor/bases/zap_logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
   @Author: orbit-w
   @File: base
   @2024 8月 周二 23:04
*/

type logStage int8

var (
	Debug logStage = 0
	Dev   logStage = 1
	Prod  logStage = 2
)

var (
	baseLogger *zap.Logger
)

const (
	FlagLogDir = "log_dir"  //日志存储路径
	FlagV      = "v"        //日志等级
	FlagStage  = "logStage" //服务环境
)

const d = 2

func NewZapLogger() *zap.Logger {
	var (
		dir string
	)

	lv := viper.GetString(FlagV)
	zlv := selectLevel(lv)

	dir = viper.GetString(FlagLogDir)
	if dir == "" {
		dir = "./transport.log"
	}

	stage := viper.GetString(FlagStage)
	switch selectStage(stage) {
	case Prod:
		return zap_logger.NewProductionLogger(dir, zlv)
	default:
		return zap_logger.NewDevelopmentLogger()
	}
}

func selectLevel(lv string) zapcore.Level {
	switch lv {
	case "Info", "info", "INFO":
		return zap.InfoLevel
	case "debug", "Debug", "DEBUG":
		return zap.DebugLevel
	case "warn", "Warn", "WARN":
		return zap.WarnLevel
	case "err", "Err", "error", "Error", "ERROR":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func selectStage(stage string) logStage {
	switch stage {
	case "Dev", "dev", "DEV", "Development", "development",
		"debug", "Debug", "DEBUG":
		return Dev
	case "Prod", "prod", "PROD", "Production", "production",
		"release", "Release", "RELEASE":
		return Prod
	default:
		return Dev
	}
}

func getBaseLogger() *zap.Logger {
	if baseLogger == nil {
		baseLogger = NewZapLogger()
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
