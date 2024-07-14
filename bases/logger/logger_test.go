package logger

import (
	"go.uber.org/zap"
	"testing"
)

/*
   @Author: orbit-w
   @File: logger_test
   @2024 1月 周二 23:19
*/

func Test_Logger(t *testing.T) {
	logger := New("./logs/test.log", zap.InfoLevel)
	logger.Info("Info record")
	logger.Error("Error record")
	Stop(logger)
}
