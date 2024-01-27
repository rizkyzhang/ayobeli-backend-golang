package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/sirupsen/logrus"
)

type loggerUtil struct {
	env    *domain.Env
	logger *logrus.Logger
}

func NewLoggerUtil(env *domain.Env) domain.LoggerUtil {
	logger := logrus.New()
	if env.AppEnv != "prod" {
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)

	}

	return &loggerUtil{env: env, logger: logger}
}

func (b *loggerUtil) EchoMiddlewareFunc() func(c echo.Context, values middleware.RequestLoggerValues) error {
	return func(c echo.Context, values middleware.RequestLoggerValues) error {
		logData := b.logger.WithFields(logrus.Fields{
			"URI":    values.URI,
			"method": values.Method,
			"status": values.Status,
		})
		if values.Status >= 300 {
			logData.Errorf("failed request with status %d", values.Status)
		} else {
			logData.Infof("success request with status %d", values.Status)
		}

		return nil
	}
}

func (b *loggerUtil) Debugf(format string, args ...interface{}) {
	b.logger.Debugf(format, args...)
}

func (b *loggerUtil) Infoln(args ...interface{}) {
	b.logger.Infoln(args)
}

func (b *loggerUtil) Infof(format string, args ...interface{}) {
	b.logger.Infof(format, args...)
}

func (b *loggerUtil) Warnf(format string, args ...interface{}) {
	b.logger.Warnf(format, args...)
}

func (b *loggerUtil) Errorf(format string, args ...interface{}) {
	b.logger.Errorf(format, args...)
}

func (b *loggerUtil) Fatalf(format string, args ...interface{}) {
	b.logger.Fatalf(format, args...)
}
