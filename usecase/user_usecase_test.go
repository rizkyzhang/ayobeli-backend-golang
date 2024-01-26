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

type UserUsecaseSuite struct {
	suite.Suite
	env          *domain.Env
	db           *sqlx.DB
	pool         *dockertest.Pool
	resource     *dockertest.Resource
	ctx          context.Context
	now          time.Time
	nowUTC       time.Time
	userRepo     domain.UserRepository
	userRepoMock *mocks.UserRepositoryMock
	user         *domain.UserModel
}

func (s *UserUsecaseSuite) BeforeTest(suiteName, testName string) {
	userPayload := &domain.UserRepositoryPayloadCreateUser{
		UID:          gofakeit.UUID(),
		FirebaseUID:  gofakeit.UUID(),
		Email:        gofakeit.Email(),
		Name:         gofakeit.Name(),
		Phone:        gofakeit.Phone(),
		ProfileImage: gofakeit.URL(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	userID, err := s.userRepo.CreateUser(userPayload)
	if err != nil {
		log.Fatal(err)
	}
	adminPayload := &domain.UserRepositoryPayloadCreateAdmin{
		UID:       gofakeit.UUID(),
		Email:     gofakeit.Email(),
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err = s.userRepo.CreateAdmin(adminPayload)
	if err != nil {
		log.Fatal(err)
	}

	user, err := s.userRepo.GetUserByUID(userPayload.UID)
	if err != nil {
		log.Fatal(err)
	}
	s.user = user
}

func (s *UserUsecaseSuite) SetupTest() {
	env := utils.LoadConfig("../.env")
	pool, resource, db := utils.SetupTestDB(env)
	s.env = env
	s.pool = pool
	s.resource = resource
	s.db = db

	ctx := context.Background()
	now := time.Now()
	userRepo := repository.NewUserRepository(s.db)
	userRepoMock := &mocks.UserRepositoryMock{}

	s.ctx = ctx
	s.now = now
	s.nowUTC = now.UTC()
	s.userRepo = userRepo
	s.userRepoMock = userRepoMock
}

func (s *UserUsecaseSuite) TearDownTest() {
	if err := s.pool.Purge(s.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestUserUsecase(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}

func (s *UserUsecaseSuite) TestUserUsecase() {
	s.Run("Get user by firebase uid should be successful", func() {
		uc := usecase.NewUserUsecase(s.env, s.userRepo)
		user, err := uc.GetUserByFirebaseUID(s.user.FirebaseUID)
		s.NoError(err)
		s.Equal(s.user, user)
	})

	s.Run("Get user by uid should be successful", func() {
		uc := usecase.NewUserUsecase(s.env, s.userRepo)
		user, err := uc.GetUserByUID(s.user.UID)
		s.NoError(err)
		s.Equal(s.user, user)
	})

	s.Run("Get admin by user id should be successful", func() {
		uc := usecase.NewUserUsecase(s.env, s.userRepo)
		admin, err := uc.GetAdminByUserID(s.user.ID)
		s.NoError(err)
		s.Equal(s.user.ID, admin.UserID)
	})
}
