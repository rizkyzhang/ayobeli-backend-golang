package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils/response_util"
)

type baseAuthController struct {
	authUsecase domain.AuthUsecase
	env         *domain.Env
	validate    *validator.Validate
}

func NewAuthController(authUsecase domain.AuthUsecase, env *domain.Env, validate *validator.Validate) domain.AuthController {
	return &baseAuthController{
		authUsecase: authUsecase,
		env:         env,
		validate:    validate,
	}
}

func (b *baseAuthController) SignUp(c echo.Context) error {
	email := c.Request().Header.Get("email")
	password := c.Request().Header.Get("password")

	err := b.validate.Var(email, "email")
	if err != nil {
		return response_util.FromBadRequestError(errors.New("invalid email")).WithEcho(c)
	}
	err = b.validate.Var(password, "min=8")
	if err != nil {
		return response_util.FromBadRequestError(errors.New("min password length is 8")).WithEcho(c)
	}

	accessToken, refreshToken, accessTokenExpirationTime, refreshTokenExpirationTime, err := b.authUsecase.SignUp(email, password)
	if err != nil {
		if err.Error() == "user already exist" {
			return response_util.FromBadRequestError(err).WithEcho(c)
		}

		return response_util.FromError(err).WithEcho(c)
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  accessTokenExpirationTime,
		HttpOnly: true,
		Secure:   b.env.AppEnv != "dev",
		Path:     "/",
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshTokenExpirationTime,
		HttpOnly: true,
		Secure:   b.env.AppEnv != "dev",
		Path:     "/",
	})

	c.SetCookie(&http.Cookie{
		Name:     "is_auth",
		Value:    "true",
		Expires:  accessTokenExpirationTime,
		HttpOnly: false,
		Path:     "/",
	})

	return response_util.FromCreated().WithEcho(c)
}

func (b *baseAuthController) SignIn(c echo.Context) error {
	email := c.Request().Header.Get("email")
	password := c.Request().Header.Get("password")

	err := b.validate.Var(email, "email")
	if err != nil {
		return response_util.FromBadRequestError(errors.New("invalid email")).WithEcho(c)
	}
	if password == "" {
		return response_util.FromBadRequestError(errors.New("password is empty")).WithEcho(c)
	}

	accessToken, refreshToken, accessTokenExpirationTime, refreshTokenExpirationTime, err := b.authUsecase.SignIn(email, password)
	if err != nil {
		if err.Error() == "user not found" || err.Error() == "wrong password" {
			return response_util.FromBadRequestError(err).WithEcho(c)
		}

		return response_util.FromError(err).WithEcho(c)
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  accessTokenExpirationTime,
		HttpOnly: true,
		Secure:   b.env.AppEnv != "dev",
		Path:     "/",
	})

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  refreshTokenExpirationTime,
		HttpOnly: true,
		Secure:   b.env.AppEnv != "dev",
		Path:     "/",
	})

	c.SetCookie(&http.Cookie{
		Name:     "is_auth",
		Value:    "true",
		Expires:  accessTokenExpirationTime,
		HttpOnly: false,
		Secure:   b.env.AppEnv != "dev",
		Path:     "/",
	})

	return response_util.FromOK().WithEcho(c)
}

func (b *baseAuthController) SignOut(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:    "access_token",
		Expires: time.Now(),
		Path:    "/",
	})

	c.SetCookie(&http.Cookie{
		Name:    "refresh_token",
		Expires: time.Now(),
		Path:    "/",
	})

	c.SetCookie(&http.Cookie{
		Name:    "is_auth",
		Expires: time.Now(),
		Path:    "/",
	})

	return response_util.FromOK().WithEcho(c)
}

func (b *baseAuthController) RefreshAccessToken(c echo.Context) error {
	cookieToken, err := c.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return response_util.FromBadRequestError(err).WithEcho(c)
		}

		return response_util.FromInternalServerError().WithEcho(c)
	}
	token := cookieToken.Value

	accessToken, accessTokenExpirationTime, err := b.authUsecase.RefreshAccessToken(token)
	if err != nil {
		return response_util.FromInternalServerError().WithEcho(c)
	}

	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  accessTokenExpirationTime,
		HttpOnly: true,
		Secure:   b.env.AppEnv != "dev",
		Path:     "/",
	})

	c.SetCookie(&http.Cookie{
		Name:     "is_auth",
		Value:    "true",
		Expires:  accessTokenExpirationTime,
		HttpOnly: false,
		Secure:   b.env.AppEnv != "dev",
		Path:     "/",
	})

	return response_util.FromOK().WithEcho(c)
}
