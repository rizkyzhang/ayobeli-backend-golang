package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	route "github.com/rizkyzhang/ayobeli-backend-golang/api/route"
	"github.com/rizkyzhang/ayobeli-backend-golang/bootstrap"
	docs "github.com/rizkyzhang/ayobeli-backend-golang/docs"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title			Ayobeli API
//	@version		0.1
//	@description	Yet another e-commerce API

//	@license.name	MIT
//	@license.url	https://opensource.org/license/MIT

//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Firebase auth access token, get it from POST /auth/access-token

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.DB
	defer app.CloseDBConnection()
	firebaseAuth := app.FirebaseAuth
	loggerUtil := utils.NewLoggerUtil(env)
	docs.SwaggerInfo.Host = env.Host

	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogError:      true,
		LogMethod:     true,
		LogStatus:     true,
		LogValuesFunc: loggerUtil.EchoMiddlewareFunc(),
	}))

	route.Setup(env, loggerUtil, db, firebaseAuth, e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(env.Port))
}
