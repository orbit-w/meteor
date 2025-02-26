package mysqldb

import (
	"strings"

	"gorm.io/gorm/logger"
)

// getLogLevel 获取GORM日志级别
func getLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Silent
	}
}
