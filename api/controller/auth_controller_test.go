package controller_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/api/controller"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain/mocks"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils/response_util"
	"github.com/stretchr/testify/suite"
)

type AuthControllerSuite struct {
	suite.Suite
	c                      echo.Context
	e                      *echo.Echo
	req                    *http.Request
	rec                    *httptest.ResponseRecorder
	ucMock                 *mocks.AuthUsecaseMock
	ct                     domain.AuthController
	badRequestRes          response_util.Response
	internalServerErrorRes response_util.Response
	email                  string
	password               string
}

func (s *AuthControllerSuite) SetupTest() {
	authUsecaseMock := &mocks.AuthUsecaseMock{}
	env := utils.LoadConfig("../../.env")
	validate := validator.New()
	ct := controller.NewAuthController(authUsecaseMock, env, validate)

	s.ct = ct
	s.ucMock = authUsecaseMock
	s.badRequestRes = response_util.Response{
		Code:   http.StatusBadRequest,
		Status: http.StatusText(http.StatusBadRequest),
	}
	s.internalServerErrorRes = response_util.Response{
		Code:   http.StatusInternalServerError,
		Status: http.StatusText(http.StatusInternalServerError),
	}
	s.email = "test@email.com"
	s.password = "test1234"
}

func (s *AuthControllerSuite) SetupSubTest() {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s.e = e
	s.req = req
	s.rec = rec
	s.c = c
}

func TestAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerSuite))
}

func (s *AuthControllerSuite) ValidateRes(rec *httptest.ResponseRecorder, expectedRes response_util.Response) {
	_res := rec.Result()
	defer _res.Body.Close()

	data, err := io.ReadAll(_res.Body)
	s.NoError(err)
	s.NotNil(data)

	var res response_util.Response
	err = json.Unmarshal(data, &res)
	s.NoError(err)
	s.Equal(expectedRes, res)
}

func (s *AuthControllerSuite) TestSignUp() {
	s.Run("Return bad request given invalid email", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "invalid email"

		s.req.Header.Set("email", "test@email")
		s.req.Header.Set("password", "test1234")

		if s.NoError(s.ct.SignUp(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return bad request given invalid password", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "min password length is 8"

		s.req.Header.Set("email", "test@email.com")
		s.req.Header.Set("password", "test")

		if s.NoError(s.ct.SignUp(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return bad request when user already exist", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "user already exist"

		s.req.Header.Set("email", s.email)
		s.req.Header.Set("password", s.password)

		s.ucMock.SignUpReturns("", "", time.Time{}, time.Time{}, errors.New("user already exist"))

		if s.NoError(s.ct.SignUp(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return internal server error when business logic failed", func() {
		expectedRes := s.internalServerErrorRes

		s.req.Header.Set("email", s.email)
		s.req.Header.Set("password", s.password)

		s.ucMock.SignUpReturns("", "", time.Time{}, time.Time{}, errors.New(""))

		if s.NoError(s.ct.SignUp(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return created", func() {
		expectedRes := response_util.Response{
			Code:   http.StatusCreated,
			Status: http.StatusText(http.StatusCreated),
		}

		s.req.Header.Set("email", s.email)
		s.req.Header.Set("password", s.password)

		s.ucMock.SignUpReturns("", "", time.Time{}, time.Time{}, nil)

		if s.NoError(s.ct.SignUp(s.c)) {
			s.ValidateRes(s.rec, expectedRes)

			cookies := s.rec.Result().Cookies()
			var cookieNames []string

			for _, cookie := range cookies {
				cookieNames = append(cookieNames, cookie.Name)
			}

			s.Contains(cookieNames, "access_token")
			s.Contains(cookieNames, "refresh_token")
			s.Contains(cookieNames, "is_auth")
		}
	})
}

func (s *AuthControllerSuite) TestSignIn() {
	s.Run("Return bad request given invalid email", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "invalid email"

		s.req.Header.Set("email", "test@email")
		s.req.Header.Set("password", "test1234")

		if s.NoError(s.ct.SignIn(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return bad request when user not found", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "user not found"

		s.req.Header.Set("email", s.email)
		s.req.Header.Set("password", s.password)

		s.ucMock.SignInReturns("", "", time.Time{}, time.Time{}, errors.New("user not found"))

		if s.NoError(s.ct.SignIn(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return bad request given wrong password", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "wrong password"

		s.req.Header.Set("email", s.email)
		s.req.Header.Set("password", s.password)

		s.ucMock.SignInReturns("", "", time.Time{}, time.Time{}, errors.New("wrong password"))

		if s.NoError(s.ct.SignIn(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return internal server error when business logic failed", func() {
		expectedRes := s.internalServerErrorRes

		s.req.Header.Set("email", s.email)
		s.req.Header.Set("password", s.password)

		s.ucMock.SignInReturns("", "", time.Time{}, time.Time{}, errors.New(""))

		if s.NoError(s.ct.SignIn(s.c)) {
			s.ValidateRes(s.rec, expectedRes)
		}
	})

	s.Run("Return OK", func() {
		expectedRes := response_util.Response{
			Code:   http.StatusOK,
			Status: http.StatusText(http.StatusOK),
		}
		s.req.Header.Set("email", s.email)
		s.req.Header.Set("password", s.password)

		s.ucMock.SignInReturns("", "", time.Time{}, time.Time{}, nil)

		if s.NoError(s.ct.SignIn(s.c)) {
			s.ValidateRes(s.rec, expectedRes)

			cookies := s.rec.Result().Cookies()
			var cookieNames []string

			for _, cookie := range cookies {
				cookieNames = append(cookieNames, cookie.Name)
			}

			s.Contains(cookieNames, "access_token")
			s.Contains(cookieNames, "refresh_token")
			s.Contains(cookieNames, "is_auth")
		}
	})
}
