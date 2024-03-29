package domain

import (
	"context"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/user_usecase_mock.go --fake-name UserUsecaseMock . UserUsecase
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/user_repository_mock.go --fake-name UserRepositoryMock . UserRepository

// Usecase
type UserUsecase interface {
	GetUserByFirebaseUID(ctx context.Context, UID string) (*UserModel, error)
	GetUserByUID(ctx context.Context, UID string) (*UserModel, error)
	GetAdminByUserID(ctx context.Context, UserID int) (*AdminModel, error)
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

type UserRepository interface {
	CreateUser(userPayload *UserRepositoryPayloadCreateUser) (int, error)
	CreateAdmin(adminPayload *UserRepositoryPayloadCreateAdmin) error
	GetUserByEmail(email string) (*UserModel, error)
	GetUserByFirebaseUID(UID string) (*UserModel, error)
	GetUserByUID(UID string) (*UserModel, error)
	GetAdminByUserID(UserID int) (*AdminModel, error)
}

type UserRepositoryPayloadCreateUser struct {
	UID          string `db:"uid" json:"uid"`
	FirebaseUID  string `db:"firebase_uid" json:"firebase_uid"`
	Email        string `db:"email" json:"email"`
	Name         string `db:"name" json:"name"`
	Phone        string `db:"phone" json:"phone"`
	ProfileImage string `db:"profile_image" json:"profile_image"`
	IsAdmin      bool   `json:"is_admin"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserRepositoryPayloadCreateAdmin struct {
	UID    string `db:"uid" json:"uid"`
	Email  string `db:"email" json:"email"`
	UserID int    `db:"user_id" json:"user_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
