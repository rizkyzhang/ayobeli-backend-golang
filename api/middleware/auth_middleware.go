package middleware

import (
	"context"
	"errors"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils/response_util"
)

type baseAuthMiddleware struct {
	authUsecase  domain.AuthUsecase
	firebaseAuth *auth.Client
}

func NewAuthMiddleware(authUsecase domain.AuthUsecase, firebaseAuth *auth.Client) domain.AuthMiddleware {
	return &baseAuthMiddleware{
		authUsecase:  authUsecase,
		firebaseAuth: firebaseAuth,
	}
}

func (b *baseAuthMiddleware) ValidateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_bearerToken := c.Request().Header.Get("authorization")
			if !strings.HasPrefix(_bearerToken, "Bearer ") {
				return response_util.FromForbiddenError(errors.New("invalid access token")).WithEcho(c)
			}

			bearerToken := strings.Split(_bearerToken, " ")[1]
			token, err := b.firebaseAuth.VerifyIDTokenAndCheckRevoked(context.Background(), bearerToken)
			if err != nil {
				return response_util.FromForbiddenError(err).WithEcho(c)
			}
			user, err := b.authUsecase.GetUserByFirebaseUID(token.UID)
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
