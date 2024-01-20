package usecase_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain/mocks"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
	"github.com/rizkyzhang/ayobeli-backend-golang/repository"
	"github.com/rizkyzhang/ayobeli-backend-golang/usecase"
	"github.com/stretchr/testify/suite"
)

type AuthUsecaseSuite struct {
	suite.Suite
	env      *domain.Env
	db       *sqlx.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
	ctx      context.Context
	now      time.Time
	nowUTC   time.Time
	authRepo domain.AuthRepository
	cartRepo domain.CartRepository
	email    string
	password string
}

func (s *AuthUsecaseSuite) SetupTest() {
	env := utils.LoadConfig("../.env")
	pool, resource, db := utils.SetupTestDB(env)
	s.env = env
	s.pool = pool
	s.resource = resource
	s.db = db

	ctx := context.Background()
	now := time.Now()
	authRepo := repository.NewAuthRepository(s.db)
	cartRepo := repository.NewCartRepository(s.db)

	s.ctx = ctx
	s.now = now
	s.nowUTC = now.UTC()
	s.authRepo = authRepo
	s.cartRepo = cartRepo
	s.email = "test@email.com"
	s.password = "test1234"
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
	s.Run("Signup should be successful", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.authRepo, authUtilMock)
		expectedFirebaseUID := gofakeit.UUID()
		authUtilMock.CreateUserReturns(expectedFirebaseUID, nil)

		err := uc.SignUp(s.email, s.password)
		s.NoError(err)

		// Validate created user
		user, err := s.authRepo.GetUserByEmail(s.email)
		s.NoError(err)
		s.NotNil(user)
		s.Equal(s.email, user.Email)
		s.Equal(expectedFirebaseUID, user.FirebaseUID)

		// Validate created cart
		cart, err := s.cartRepo.GetCartByUserID(user.ID)
		s.NoError(err)
		s.NotNil(cart)
	})

	s.Run("Signup should return an error if user already exist", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.authRepo, authUtilMock)

		err := uc.SignUp(s.email, s.password)
		s.Error(err)
	})

	s.Run("Get access token should be successful", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.authRepo, authUtilMock)
		expectedAccessToken := gofakeit.UUID()
		authUtilMock.GetAccessTokenReturns(expectedAccessToken, nil)

		accessToken, err := uc.GetAccessToken(s.email, s.password)
		s.NoError(err)
		s.Equal(expectedAccessToken, accessToken)
	})

	s.Run("Get access token should return an error if user not found", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.authRepo, authUtilMock)

		_, err := uc.GetAccessToken("notfound@email.com", s.password)
		s.Error(err)
	})

	s.Run("Get access token should return an error if result is empty which indicate invalid password", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		authUtilMock.GetAccessTokenReturns("", nil)
		uc := usecase.NewAuthUsecase(s.env, s.authRepo, authUtilMock)

		_, err := uc.GetAccessToken(s.email, "invalid")
		s.Error(err)
	})
}
