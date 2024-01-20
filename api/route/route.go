package route

import (
	"firebase.google.com/go/v4/auth"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/rizkyzhang/ayobeli-backend-golang/api/middleware"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/repository"
	"github.com/rizkyzhang/ayobeli-backend-golang/usecase"
)

func Setup(env *domain.Env, db *sqlx.DB, firebaseAuth *auth.Client, e *echo.Echo) {
	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(env, authRepo, firebaseAuth)
	authMiddleware := middleware.NewAuthMiddleware(authUsecase, firebaseAuth)
	validate := validator.New()

	rootGroup := e.Group("/api")

	NewAuthRouter(env, db, rootGroup, authUsecase, authMiddleware, validate)
}
