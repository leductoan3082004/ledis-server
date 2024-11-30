package logging

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	logger     *logrus.Logger
	loggerOnce sync.Once
)

func initLogger() {
	logger = logrus.New()

	logger.SetOutput(os.Stdout)

	logger.SetFormatter(
		&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			DisableQuote:    true,
		},
	)

	logger.SetLevel(logrus.DebugLevel)
}

func GetLogger() *logrus.Logger {
	loggerOnce.Do(initLogger)
	return logger
}
