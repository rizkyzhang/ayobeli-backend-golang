package route

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/api/controller"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

func NewAuthRouter(env *domain.Env, db *sqlx.DB, rootGroup *echo.Group, authUsecase domain.AuthUsecase, authMiddleware domain.AuthMiddleware, validate *validator.Validate) {
	ct := controller.NewAuthController(authUsecase, env, validate)

	publicGroup := rootGroup.Group("/v1/auth")
	privateGroup := rootGroup.Group("/v1/auth")
	privateGroup.Use(authMiddleware.ValidateUser())

	publicGroup.POST("/signup", ct.SignUp)
	publicGroup.GET("/token", ct.GetAccessToken)
}
