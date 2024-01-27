package middleware

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils/response_util"
)

type baseAuthMiddleware struct {
	userUsecase domain.UserUsecase
	authUtil    domain.AuthUtil
}

func NewAuthMiddleware(userUsecase domain.UserUsecase, authUtil domain.AuthUtil) domain.AuthMiddleware {
	return &baseAuthMiddleware{
		userUsecase: userUsecase,
		authUtil:    authUtil,
	}
}

func (b *baseAuthMiddleware) ValidateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerToken := c.Request().Header.Get("authorization")
			if !strings.HasPrefix(bearerToken, "Bearer ") {
				return response_util.FromForbiddenError(errors.New("invalid access token")).WithEcho(c)
			}

			token := strings.Split(bearerToken, " ")[1]
			firebaseUID, err := b.authUtil.VerifyToken(token)
			if err != nil {
				return response_util.FromForbiddenError(err).WithEcho(c)
			}
			user, err := b.userUsecase.GetUserByFirebaseUID(c.Request().Context(), firebaseUID)
			if user == nil {
				return response_util.FromForbiddenError(errors.New("access denied")).WithEcho(c)
			}
			if err != nil {
				return response_util.FromInternalServerError().WithEcho(c)
			}
			c.Set("user", user)

			return next(c)
		}
	}
}

func (b *baseAuthMiddleware) ValidateAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*domain.UserModel)
			if !ok || user == nil {
				return response_util.FromForbiddenError(errors.New("access denied")).WithEcho(c)
			}

			admin, err := b.userUsecase.GetAdminByUserID(c.Request().Context(), user.ID)
			if admin == nil {
				return response_util.FromForbiddenError(errors.New("access denied")).WithEcho(c)
			}
			if err != nil {
				return response_util.FromInternalServerError().WithEcho(c)
			}
			c.Set("admin", admin)

			return next(c)
		}
	}
}
