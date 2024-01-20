package middleware

import (
	"context"
	"errors"
	"net/http"

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

type AuthHeader struct {
	Authorization string `header:"authorization"`
}

func (b *baseAuthMiddleware) ValidateUser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var authHeader AuthHeader

			err := c.Bind(&authHeader)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					return response_util.FromForbiddenError(errors.New("access token not found")).WithEcho(c)
				}

				return response_util.FromForbiddenError(errors.New("invalid access token")).WithEcho(c)
			}

			token, err := b.firebaseAuth.VerifyIDTokenAndCheckRevoked(context.Background(), authHeader.Authorization)
			if err != nil {
				return response_util.FromForbiddenError(err).WithEcho(c)
			}
			user, err := b.authUsecase.GetUserByUID(token.UID)
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
