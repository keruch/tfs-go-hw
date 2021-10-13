package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger = logrus.Logger

func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	return logger
}
