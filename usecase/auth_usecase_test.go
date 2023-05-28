package usecase_test

import (
	"context"
	"log"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
	"github.com/rizkyzhang/ayobeli-backend/repository"
	"github.com/rizkyzhang/ayobeli-backend/usecase"
	"github.com/stretchr/testify/suite"
)

type AuthUsecaseSuite struct {
	suite.Suite
	db           *sqlx.DB
	pool         *dockertest.Pool
	resource     *dockertest.Resource
	ctx          context.Context
	now          time.Time
	nowUTC       time.Time
	authRepo     domain.AuthRepository
	cartRepo     domain.CartRepository
	hashUtil     domain.HashUtil
	jwtUtil      domain.JWTUtil
	email        string
	password     string
	accessToken  string
	refreshToken string
}

func (s *AuthUsecaseSuite) SetupTest() {
	env := utils.LoadConfig("../.env")
	pool, resource, db := utils.SetupTestDB(env)

	s.pool = pool
	s.resource = resource
	s.db = db

	ctx := context.Background()
	now := time.Now()
	authRepo := repository.NewAuthRepository(s.db)
	cartRepo := repository.NewCartRepository(s.db)
	hashUtil := utils.NewHashUtil()
	jwtUtil := utils.NewJWTUtil([]byte(env.AccessTokenSecret), []byte(env.RefreshTokenSecret), env.AccessTokenExpiryHour, env.RefreshTokenExpiryHour)

	s.ctx = ctx
	s.now = now
	s.nowUTC = now.UTC()
	s.authRepo = authRepo
	s.cartRepo = cartRepo
	s.hashUtil = hashUtil
	s.jwtUtil = jwtUtil
	s.email = "test@email.com"
	s.password = "test123"
}

func (s *AuthUsecaseSuite) TearDownTest() {
	if err := s.pool.Purge(s.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestAuthUsecase(t *testing.T) {
	suite.Run(t, new(AuthUsecaseSuite))
}

func (s *AuthUsecaseSuite) TestAuthUsecase() {
	s.Run("Signup user", func() {
		uc := usecase.NewAuthUsecase(s.authRepo, s.hashUtil, s.jwtUtil)

		accessToken, refreshToken, accessTokenExpirationTime, refreshTokenExpirationTime, err := uc.SignUp(s.email, s.password)
		s.NoError(err)
		s.NotEmpty(accessTokenExpirationTime)
		s.NotEmpty(refreshTokenExpirationTime)

		// Validate generated token
		_, err = s.jwtUtil.ParseUserUID(accessToken, true)
		s.NoError(err)
		_, err = s.jwtUtil.ParseUserUID(refreshToken, false)
		s.NoError(err)

		// Validate created user
		user, err := s.authRepo.GetUserByEmail(s.email)
		s.NoError(err)
		s.NotNil(user)
		s.Equal(s.email, user.Email)
		isValidPassword := s.hashUtil.ValidatePassword(s.password, user.Password)
		s.True(isValidPassword)

		// Validate created cart
		cart, err := s.cartRepo.GetCartByUserID(user.ID)
		s.NoError(err)
		s.NotNil(cart)
	})

	s.Run("Signup user return an error if user already exist", func() {
		uc := usecase.NewAuthUsecase(s.authRepo, s.hashUtil, s.jwtUtil)

		_, _, _, _, err := uc.SignUp(s.email, s.password)
		s.Error(err)
	})

	s.Run("Signin user", func() {
		uc := usecase.NewAuthUsecase(s.authRepo, s.hashUtil, s.jwtUtil)

		accessToken, refreshToken, accessTokenExpirationTime, refreshTokenExpirationTime, err := uc.SignIn(s.email, s.password)
		s.NoError(err)
		s.NotEmpty(accessTokenExpirationTime)
		s.NotEmpty(refreshTokenExpirationTime)

		s.accessToken = accessToken
		s.refreshToken = refreshToken

		// Validate generated token
		_, err = s.jwtUtil.ParseUserUID(accessToken, true)
		s.NoError(err)
		_, err = s.jwtUtil.ParseUserUID(refreshToken, false)
		s.NoError(err)
	})

	s.Run("Signin return an error if user not found", func() {
		uc := usecase.NewAuthUsecase(s.authRepo, s.hashUtil, s.jwtUtil)

		_, _, _, _, err := uc.SignIn("notfound@email.com", s.password)
		s.Error(err)
	})

	s.Run("Signin return an error if password is not valid", func() {
		uc := usecase.NewAuthUsecase(s.authRepo, s.hashUtil, s.jwtUtil)

		_, _, _, _, err := uc.SignIn(s.email, "test321")
		s.Error(err)
	})

	s.Run("Refresh access token", func() {
		uc := usecase.NewAuthUsecase(s.authRepo, s.hashUtil, s.jwtUtil)

		accessToken, accessTokenExpirationTime, err := uc.RefreshAccessToken(s.refreshToken)
		s.NoError(err)
		s.NotEmpty(accessTokenExpirationTime)

		// Validate generated token
		_, err = s.jwtUtil.ParseUserUID(accessToken, true)
		s.NoError(err)
	})

	s.Run("Refresh access token return an error given invalid refresh token", func() {
		uc := usecase.NewAuthUsecase(s.authRepo, s.hashUtil, s.jwtUtil)

		_, _, err := uc.RefreshAccessToken(s.accessToken)
		s.Error(err)
	})
}
