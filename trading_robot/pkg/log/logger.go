package log

import (
		"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger

func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	return logger
}
