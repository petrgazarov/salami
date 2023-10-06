package logger

import (
	"os"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	if os.Getenv("DEBUG") != "true" {
		return
	}

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = zapLogger.Sugar()
}

func Log(message string) {
	if logger == nil {
		return
	}

	defer logger.Sync()
	logger.Info(message)
}
