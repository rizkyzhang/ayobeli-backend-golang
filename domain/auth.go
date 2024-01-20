package domain

import (
	"time"

	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/auth_usecase_mock.go --fake-name AuthUsecaseMock . AuthUsecase

// Controller
type AuthMiddleware interface {
	ValidateUser() echo.MiddlewareFunc
	ValidateAdmin() echo.MiddlewareFunc
}

type AuthController interface {
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
	GetUserByUID(UID string) (*UserModel, error)
	GetAdminByUserID(UserID int) (*AdminModel, error)
}

// Repository
type UserModel struct {
	ID           int    `db:"id" json:"id"`
	UID          string `db:"uid" json:"uid"`
	FirebaseUID  string `db:"firebase_uid" json:"firebase_uid"`
	Email        string `db:"email" json:"email"`
	Name         string `db:"name" json:"name"`
	Phone        string `db:"phone" json:"phone"`
	ProfileImage string `db:"profile_image" json:"profile_image"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type AdminModel struct {
	ID     int    `db:"id" json:"id"`
	UID    string `db:"uid" json:"uid"`
	Email  string `db:"email" json:"email"`
	UserID int    `db:"user_id" json:"user_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type AuthRepository interface {
	CreateUser(userPayload *AuthRepositoryPayloadCreateUser) (int, error)
	GetUserByEmail(email string) (*UserModel, error)
	GetUserByUID(UID string) (*UserModel, error)
	GetAdminByUserID(UserID int) (*AdminModel, error)
}

type AuthRepositoryPayloadCreateUser struct {
	UID          string `db:"uid" json:"uid"`
	FirebaseUID  string `db:"firebase_uid" json:"firebase_uid"`
	Email        string `db:"email" json:"email"`
	Name         string `db:"name" json:"name"`
	Phone        string `db:"phone" json:"phone"`
	ProfileImage string `db:"profile_image" json:"profile_image"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
