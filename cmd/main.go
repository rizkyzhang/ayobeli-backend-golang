package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rizkyzhang/ayobeli-backend/bootstrap"
	"github.com/sirupsen/logrus"
)

func main() {
	app := bootstrap.App()
	env := app.Env

	if env.AppEnv != "prod" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogError:  true,
		LogMethod: true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			logData := logrus.WithFields(logrus.Fields{
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
		},
	}))

	e.Logger.Fatal(e.Start(env.ServerAddress))
}
