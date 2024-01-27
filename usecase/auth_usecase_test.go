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
	userRepo domain.UserRepository
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
	userRepo := repository.NewUserRepository(s.db)
	cartRepo := repository.NewCartRepository(s.db)

	s.ctx = ctx
	s.now = now
	s.nowUTC = now.UTC()
	s.userRepo = userRepo
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
	s.Run("Signup as admin should be successful", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.userRepo, authUtilMock)
		expectedFirebaseUID := gofakeit.UUID()
		authUtilMock.CreateUserReturns(expectedFirebaseUID, nil)

		err := uc.SignUp(s.ctx, s.email, s.password, true)
		s.NoError(err)

		// Validate created user
		user, err := s.userRepo.GetUserByEmail(s.email)
		s.NoError(err)
		s.NotNil(user)
		s.Equal(s.email, user.Email)
		s.Equal(expectedFirebaseUID, user.FirebaseUID)

		admin, err := s.userRepo.GetAdminByUserID(user.ID)
		s.NoError(err)
		s.Equal(user.ID, admin.UserID)

		// Validate created cart
		cart, err := s.cartRepo.GetCartByUserID(user.ID)
		s.NoError(err)
		s.NotNil(cart)
	})

	s.Run("Signup as non-admin should be successful", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.userRepo, authUtilMock)
		expectedFirebaseUID := gofakeit.UUID()
		authUtilMock.CreateUserReturns(expectedFirebaseUID, nil)

		email := gofakeit.Email()
		err := uc.SignUp(s.ctx, email, s.password, false)
		s.NoError(err)

		// Validate created user
		user, err := s.userRepo.GetUserByEmail(email)
		s.NoError(err)
		s.NotNil(user)
		s.Equal(email, user.Email)
		s.Equal(expectedFirebaseUID, user.FirebaseUID)

		// Validate created cart
		cart, err := s.cartRepo.GetCartByUserID(user.ID)
		s.NoError(err)
		s.NotNil(cart)
	})

	s.Run("Signup should return an error if user already exist", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.userRepo, authUtilMock)

		err := uc.SignUp(s.ctx, s.email, s.password, false)
		s.Error(err)
	})

	s.Run("Get access token should be successful", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.userRepo, authUtilMock)
		expectedAccessToken := gofakeit.UUID()
		authUtilMock.GetAccessTokenReturns(expectedAccessToken, nil)

		accessToken, err := uc.GetAccessToken(s.ctx, s.email, s.password)
		s.NoError(err)
		s.Equal(expectedAccessToken, accessToken)
	})

	s.Run("Get access token should return an error if user not found", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		uc := usecase.NewAuthUsecase(s.env, s.userRepo, authUtilMock)

		_, err := uc.GetAccessToken(s.ctx, "notfound@email.com", s.password)
		s.Error(err)
	})

	s.Run("Get access token should return an error if result is empty which indicate invalid password", func() {
		authUtilMock := &mocks.AuthUtilMock{}
		authUtilMock.GetAccessTokenReturns("", nil)
		uc := usecase.NewAuthUsecase(s.env, s.userRepo, authUtilMock)

		_, err := uc.GetAccessToken(s.ctx, s.email, "invalid")
		s.Error(err)
	})
}
