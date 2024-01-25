package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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
	ucMock                 *mocks.AuthUsecaseMock
	ct                     domain.AuthController
	createdRes             response_util.Response
	badRequestRes          response_util.Response
	notFoundRes            response_util.Response
	internalServerErrorRes response_util.Response
	reqHelper              func(body io.Reader) (echo.Context, *httptest.ResponseRecorder)
}

func (s *AuthControllerSuite) SetupTest() {
	env := utils.LoadConfig("../../.env")
	validate := validator.New()
	authUsecaseMock := &mocks.AuthUsecaseMock{}
	ct := controller.NewAuthController(authUsecaseMock, env, validate)

	s.ct = ct
	s.ucMock = authUsecaseMock
	s.createdRes = response_util.Response{
		Code:   http.StatusCreated,
		Status: http.StatusText(http.StatusCreated),
	}
	s.badRequestRes = response_util.Response{
		Code:   http.StatusBadRequest,
		Status: http.StatusText(http.StatusBadRequest),
	}
	s.notFoundRes = response_util.Response{
		Code:   http.StatusNotFound,
		Status: http.StatusText(http.StatusNotFound),
	}
	s.internalServerErrorRes = response_util.Response{
		Code:   http.StatusInternalServerError,
		Status: http.StatusText(http.StatusInternalServerError),
	}
	s.reqHelper = func(body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
		req := httptest.NewRequest(http.MethodPost, "/", body)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		rec := httptest.NewRecorder()
		e := echo.New()
		c := e.NewContext(req, rec)

		return c, rec
	}
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
	s.Run("Signup as admin should return created if successful", func() {
		expectedRes := s.createdRes

		reqBody := &domain.AuthControllerPayloadSignUp{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
			IsAdmin:  true,
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.SignUpReturns(nil)
		if s.NoError(s.ct.SignUp(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Signup as non-admin should return created if successful", func() {
		expectedRes := s.createdRes

		reqBody := &domain.AuthControllerPayloadSignUp{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
			IsAdmin:  false,
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.SignUpReturns(nil)
		if s.NoError(s.ct.SignUp(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Signup should return bad request error given invalid email", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "invalid email"

		reqBody := &domain.AuthControllerPayloadSignUp{
			Email:    "invalid",
			Password: gofakeit.Password(true, true, true, true, false, 8),
			IsAdmin:  false,
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		if s.NoError(s.ct.SignUp(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Signup should return bad request error given invalid password", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "min password length is 8"

		reqBody := &domain.AuthControllerPayloadSignUp{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 5),
			IsAdmin:  false,
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.SignUpReturns(nil)
		if s.NoError(s.ct.SignUp(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Signup should return bad request error if user already exist", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "user already exist"

		reqBody := &domain.AuthControllerPayloadSignUp{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
			IsAdmin:  false,
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.SignUpReturns(errors.New("user already exist"))
		if s.NoError(s.ct.SignUp(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Signup should return internal server error if business logic failed", func() {
		expectedRes := s.internalServerErrorRes

		reqBody := &domain.AuthControllerPayloadSignUp{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
			IsAdmin:  false,
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.SignUpReturns(errors.New(""))
		if s.NoError(s.ct.SignUp(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})
}

func (s *AuthControllerSuite) TestSignIn() {
	s.Run("Get access token should return OK if successful", func() {
		expectedData := gofakeit.UUID()
		expectedRes := response_util.Response{
			Code:   http.StatusOK,
			Status: http.StatusText(http.StatusOK),
			Data:   expectedData,
		}

		reqBody := &domain.AuthControllerPayloadGetAccessToken{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.GetAccessTokenReturns(expectedData, nil)
		if s.NoError(s.ct.GetAccessToken(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Get access token should return bad request given invalid email", func() {
		expectedRes := s.badRequestRes
		expectedRes.Error = "invalid email"

		reqBody := &domain.AuthControllerPayloadGetAccessToken{
			Email:    "invalid",
			Password: gofakeit.Password(true, true, true, true, false, 8),
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		if s.NoError(s.ct.GetAccessToken(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Get access token should return not found error if user not found", func() {
		expectedRes := s.notFoundRes
		expectedRes.Error = "user not found"

		reqBody := &domain.AuthControllerPayloadGetAccessToken{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.GetAccessTokenReturns("", errors.New("user not found"))
		if s.NoError(s.ct.GetAccessToken(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})

	s.Run("Get access token should return internal server error if business logic failed", func() {
		expectedRes := s.internalServerErrorRes

		reqBody := &domain.AuthControllerPayloadGetAccessToken{
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 8),
		}
		reqBytes, err := json.Marshal(reqBody)
		s.NoError(err)
		c, rec := s.reqHelper(bytes.NewBuffer(reqBytes))

		s.ucMock.GetAccessTokenReturns("", errors.New(""))
		if s.NoError(s.ct.GetAccessToken(c)) {
			s.ValidateRes(rec, expectedRes)
		}
	})
}
