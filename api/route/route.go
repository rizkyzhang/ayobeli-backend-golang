package route

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/rizkyzhang/ayobeli-backend/api/middleware"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
	"github.com/rizkyzhang/ayobeli-backend/repository"
	"github.com/rizkyzhang/ayobeli-backend/usecase"
)

func Setup(env *domain.Env, db *sqlx.DB, e *echo.Echo) {
	hashUtil := utils.NewHashUtil()
	jwtUtil := utils.NewJWTUtil([]byte(env.AccessTokenSecret), []byte(env.RefreshTokenSecret), env.AccessTokenExpiryHour, env.RefreshTokenExpiryHour)

	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo, hashUtil, jwtUtil)
	authMiddleware := middleware.NewAuthMiddleware(authUsecase, jwtUtil)
	validate := validator.New()

	rootGroup := e.Group("/api")

	NewAuthRouter(env, db, rootGroup, authUsecase, authMiddleware, validate)
}
