package logger

import (
	"os"

	"go.uber.org/zap"
)

type salamiLogger struct {
	verbose  bool
	instance *zap.SugaredLogger
}

var logger *salamiLogger

func InitializeLogger(verbose bool) {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeCaller = nil
	zapConfig.EncoderConfig.LevelKey = ""

	zapLogger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	logger = &salamiLogger{
		verbose:  verbose,
		instance: zapLogger.Sugar(),
	}
}

// Log logs the message always
func Log(message string) {
	if logger == nil {
		return
	}

	defer logger.instance.Sync()
	logger.instance.Info(message)
}

// Verbose logs the message if the verbose flag is set to true
func Verbose(message string) {
	if logger == nil {
		return
	}
	if !logger.verbose {
		return
	}

	defer logger.instance.Sync()
	logger.instance.Info(message)
}

// Debug logs the message if the DEBUG environment variable is set to true
func Debug(message string) {
	if logger == nil {
		return
	}
	if os.Getenv("DEBUG") != "true" {
		return
	}

	defer logger.instance.Sync()
	logger.instance.Debug(message)
}
