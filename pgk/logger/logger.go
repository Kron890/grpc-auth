package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// запуск логгера
func Init() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04 02-01-2006",
	})

	return log
}
