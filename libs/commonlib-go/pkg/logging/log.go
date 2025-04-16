package logging

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
)

func NewLogHelper(
	cfg *config.AppConfig,
) *LogHelper {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})
	logger.SetOutput(os.Stdout)
	setLoggingLevel(logger, cfg.Logging.Level)

	log.SetOutput(logger.Writer())
	log.SetFlags(0)

	return &LogHelper{Logger: logger}
}

type LogHelper struct {
	Logger *logrus.Logger
}

func setLoggingLevel(logger *logrus.Logger, level string) {
	switch level {
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
		logrus.SetLevel(logrus.FatalLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
		logrus.SetLevel(logrus.InfoLevel)
	}
}
