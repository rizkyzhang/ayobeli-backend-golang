package controller

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils/response_util"
)

type baseAuthController struct {
	env         *domain.Env
	loggerUtil  domain.LoggerUtil
	authUsecase domain.AuthUsecase
	validate    *validator.Validate
}

func NewAuthController(env *domain.Env, loggerUtil domain.LoggerUtil, authUsecase domain.AuthUsecase, validate *validator.Validate) domain.AuthController {
	return &baseAuthController{
		env:         env,
		loggerUtil:  loggerUtil,
		authUsecase: authUsecase,
		validate:    validate,
	}
}

// SignUp godoc
//
//	@Summary	Create user
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		credential body	domain.AuthControllerPayloadSignUp true	"email and password"
//	@Success	201
//	@Failure	400	"validation error | user already exist"
//	@Failure	500	"Internal Server Error"
//	@Router		/auth/signup [post]
func (b *baseAuthController) SignUp(c echo.Context) error {
	var payload domain.AuthControllerPayloadSignUp
	err := c.Bind(&payload)
	if err != nil {
		return response_util.FromBindingError(err).WithEcho(c)
	}
	err = b.validate.Struct(&payload)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return response_util.FromValidationErrors(validationErrors).WithEcho(c)
		}
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

// GetAccessToken godoc
//
//	@Summary	Get access token
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		credential body domain.AuthControllerPayloadGetAccessToken	true "email and password"
//	@Success	200
//	@Failure	400	"validation error"
//	@Failure	404	"user not found"
//	@Failure	500	"Internal Server Error"
//	@Router		/auth/access-token [post]
func (b *baseAuthController) GetAccessToken(c echo.Context) error {
	var payload domain.AuthControllerPayloadGetAccessToken
	err := c.Bind(&payload)
	if err != nil {
		return response_util.FromBindingError(err).WithEcho(c)
	}
	err = b.validate.Struct(&payload)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return response_util.FromValidationErrors(validationErrors).WithEcho(c)
		}
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
