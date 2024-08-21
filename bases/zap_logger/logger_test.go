package zap_logger

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

/*
   @Author: orbit-w
   @File: logger_test
   @2024 1月 周二 23:19
*/

func Test_Logger(t *testing.T) {
	logger := NewLogger("./logs/test.log", zap.InfoLevel)
	logger.Info("Info record")
	logger.Error("Error record")
	Stop(logger)
}

func Test_LoggerV2(t *testing.T) {
	logger := NewDevelopmentLogger()
	logger.Info("Is Test", zap.String("Name", "Test"))
	logger.Error("Is Test", zap.String("Name", "Test"))
	logger.Warn("Is Test", zap.String("Name", "Test"))
	//l.DPanic("Is Test", zap.String("Name", "Test"))

	fmt.Println("TestInitLogger")
	Stop(logger)
}
