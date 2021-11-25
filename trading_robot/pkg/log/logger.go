package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger

func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: false,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
	return logger
}
