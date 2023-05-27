package domain

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Controller
type AuthMiddleware interface {
	ValidateUser() echo.MiddlewareFunc
	ValidateAdmin() echo.MiddlewareFunc
}

type AuthController interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
	SignOut(c echo.Context) error
	RefreshToken(c echo.Context) error
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
	SignUp(email, password string) (accessToken, refreshToken string, accessTokenExpirationTime, refreshTokenExpirationTime time.Time, err error)
	SignIn(email, password string) (accessToken, refreshToken string, accessTokenExpirationTime, refreshTokenExpirationTime time.Time, err error)
	RefreshToken(refreshToken string) (string, time.Time, error)
	GetUserByUID(UID string) (*UserModel, error)
	GetAdminByUserID(UserID uint64) (*AdminModel, error)
}

// Repository
type UserModel struct {
	ID           uint64 `db:"id" json:"id"`
	UID          string `db:"uid" json:"uid"`
	Email        string `db:"email" json:"email"`
	Password     string `db:"password" json:"password"`
	Name         string `db:"name" json:"name"`
	Phone        string `db:"phone" json:"phone"`
	ProfileImage string `db:"profile_image" json:"profile_image"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type AdminModel struct {
	ID     uint64 `db:"id" json:"id"`
	UID    string `db:"uid" json:"uid"`
	Email  string `db:"email" json:"email"`
	UserID uint64 `db:"user_id" json:"user_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type AuthRepository interface {
	CreateUser(userPayload *AuthRepositoryPayloadCreateUser) (uint64, error)
	GetUserByEmail(email string) (*UserModel, error)
	GetUserByUID(UID string) (*UserModel, error)
	GetAdminByUserID(UserID uint64) (*AdminModel, error)
}

type AuthRepositoryPayloadCreateUser struct {
	UID          string `db:"uid" json:"uid"`
	Email        string `db:"email" json:"email"`
	Password     string `db:"password" json:"password"`
	Name         string `db:"name" json:"name"`
	Phone        string `db:"phone" json:"phone"`
	ProfileImage string `db:"profile_image" json:"profile_image"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
