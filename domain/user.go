package domain

import "time"

// Usecase
type UserUsecase interface {
	GetUserByFirebaseUID(UID string) (*UserModel, error)
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

type UserRepository interface {
	CreateUser(userPayload *UserRepositoryPayloadCreateUser) (int, error)
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

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
