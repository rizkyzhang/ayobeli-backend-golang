package route

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/api/controller"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

func NewAuthRouter(env *domain.Env, loggerUtil domain.LoggerUtil, rootGroup *echo.Group, authUsecase domain.AuthUsecase, authMiddleware domain.AuthMiddleware, validate *validator.Validate) {
	ct := controller.NewAuthController(env, loggerUtil, authUsecase, validate)

	publicGroup := rootGroup.Group("/v1/auth")
	privateGroup := rootGroup.Group("/v1/auth")
	privateGroup.Use(authMiddleware.ValidateUser())

	publicGroup.POST("/signup", ct.SignUp)
	publicGroup.POST("/token", ct.GetAccessToken)
}
