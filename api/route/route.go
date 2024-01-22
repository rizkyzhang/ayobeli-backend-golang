package route

import (
	"firebase.google.com/go/v4/auth"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"

	"github.com/rizkyzhang/ayobeli-backend-golang/api/middleware"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/repository"
	"github.com/rizkyzhang/ayobeli-backend-golang/usecase"
)

func Setup(env *domain.Env, db *sqlx.DB, firebaseAuth *auth.Client, e *echo.Echo) {
	authUtil := utils.NewAuthUtil(env, firebaseAuth)
	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(env, userRepo, authUtil)
	userUsecase := usecase.NewUserUsecase(env, userRepo)
	authMiddleware := middleware.NewAuthMiddleware(userUsecase, authUtil)
	validate := validator.New()

	rootGroup := e.Group("/api")

	NewAuthRouter(env, rootGroup, authUsecase, authMiddleware, validate)
}
