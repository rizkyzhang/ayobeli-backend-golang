package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils/response_util"
)

type baseAuthMiddleware struct {
	authUsecase domain.AuthUsecase
	jwtUtil     domain.JWTUtil
}

func NewAuthMiddleware(authUsecase domain.AuthUsecase, jwtUtil domain.JWTUtil) domain.AuthMiddleware {
	return &baseAuthMiddleware{
		authUsecase: authUsecase,
		jwtUtil:     jwtUtil,
	}
}

func (b *baseAuthMiddleware) ValidateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookieToken, err := c.Cookie("access_token")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return response_util.FromForbiddenError(errors.New("access token not found")).WithEcho(c)
				}

				return response_util.FromForbiddenError(errors.New("invalid access token")).WithEcho(c)
			}
			accessToken := cookieToken.Value

			userUID, err := b.jwtUtil.ParseUserUID(accessToken, true)
			if err != nil {
				return response_util.FromForbiddenError(err).WithEcho(c)
			}
			user, err := b.authUsecase.GetUserByUID(userUID)
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

			admin, err := b.authUsecase.GetAdminByUserID(user.ID)
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
