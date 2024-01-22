package domain

import (
	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/auth_usecase_mock.go --fake-name AuthUsecaseMock . AuthUsecase

// Controller
type AuthMiddleware interface {
	ValidateUser() echo.MiddlewareFunc
	ValidateAdmin() echo.MiddlewareFunc
}

type AuthController interface {
	GetAccessToken(c echo.Context) error
	SignUp(c echo.Context) error
}

type AuthControllerPayloadSignUp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthControllerPayloadSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Usecase
type AuthUsecase interface {
	SignUp(email, password string) error
	GetAccessToken(email, password string) (string, error)
}
