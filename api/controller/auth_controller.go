package controller

import (
	"errors"

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
	var payload domain.AuthControllerPayloadSignUp
	err := c.Bind(&payload)
	if err != nil {
		return response_util.FromBindingError(err).WithEcho(c)
	}
	err = b.validate.Var(payload.Email, "email")
	if err != nil {
		return response_util.FromBadRequestError(errors.New("invalid email")).WithEcho(c)
	}
	err = b.validate.Var(payload.Password, "min=8")
	if err != nil {
		return response_util.FromBadRequestError(errors.New("min password length is 8")).WithEcho(c)
	}

	err = b.authUsecase.SignUp(payload.Email, payload.Password, payload.IsAdmin)
	if err != nil {
		if err.Error() == "user already exist" {
			return response_util.FromBadRequestError(err).WithEcho(c)
		}

		return response_util.FromError(err).WithEcho(c)
	}

	return response_util.FromCreated().WithEcho(c)
}

func (b *baseAuthController) GetAccessToken(c echo.Context) error {
	var payload domain.AuthControllerPayloadGetAccessToken
	err := c.Bind(&payload)
	if err != nil {
		return response_util.FromBindingError(err).WithEcho(c)
	}
	err = b.validate.Var(payload.Email, "email")
	if err != nil {
		return response_util.FromBadRequestError(errors.New("invalid email")).WithEcho(c)
	}
	err = b.validate.Var(payload.Password, "min=8")
	if err != nil {
		return response_util.FromBadRequestError(errors.New("min password length is 8")).WithEcho(c)
	}

	accessToken, err := b.authUsecase.GetAccessToken(payload.Email, payload.Password)
	if err != nil {
		if err.Error() == "user not found" {
			return response_util.FromNotFoundError(err).WithEcho(c)
		}

		return response_util.FromError(err).WithEcho(c)
	}

	return response_util.FromData(accessToken).WithEcho(c)
}
