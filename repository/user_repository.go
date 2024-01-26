package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
)

type baseUserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &baseUserRepository{db: db}
}

func (b *baseUserRepository) CreateUser(userPayload *domain.UserRepositoryPayloadCreateUser) (int, error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		tx.Rollback()
	}()

	query, args, err := tx.BindNamed(`
	INSERT INTO users (uid, firebase_uid, email, name, phone, profile_image, created_at, updated_at)
	VALUES (:uid, :firebase_uid, :email, :name, :phone, :profile_image, :created_at, :updated_at)
	RETURNING id;
	`, userPayload)
	if err != nil {
		return 0, err
	}
	var userID int
	err = tx.Get(&userID, query, args...)
	if err != nil {
		return 0, err
	}

	metadata := utils.GenerateMetadata()
	if userPayload.IsAdmin {
		_, err := tx.Exec(`
		INSERT INTO admins (uid, email, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5);
		`, metadata.UID(), userPayload.Email, userID, metadata.CreatedAt, metadata.UpdatedAt)
		if err != nil {
			return 0, err
		}
	}

	metadata = utils.GenerateMetadata()
	cart := domain.CartModel{
		UID:              metadata.UID(),
		Quantity:         0,
		TotalPrice:       "Rp 0",
		TotalPriceValue:  0,
		TotalWeight:      "Rp 0",
		TotalWeightValue: 0,
		UserID:           userID,
		CreatedAt:        metadata.CreatedAt,
		UpdatedAt:        metadata.UpdatedAt,
	}
	_, err = tx.NamedExec(`
	INSERT INTO carts (uid, quantity, total_price, total_price_value, total_weight, total_weight_value, user_id, created_at, updated_at)
	VALUES (:uid, :quantity, :total_price, :total_price_value, :total_weight, :total_weight_value, :user_id, :created_at, :updated_at)
	`, cart)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (b *baseUserRepository) CreateAdmin(adminPayload *domain.UserRepositoryPayloadCreateAdmin) error {
	_, err := b.db.Exec(`
		INSERT INTO admins (uid, email, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5);
		`, adminPayload.UID, adminPayload.Email, adminPayload.UserID, adminPayload.CreatedAt, adminPayload.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseUserRepository) GetUserByEmail(email string) (*domain.UserModel, error) {
	var user domain.UserModel

	err := b.db.Get(&user, "SELECT * FROM users WHERE email = $1;", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (b *baseUserRepository) GetUserByFirebaseUID(UID string) (*domain.UserModel, error) {
	var user domain.UserModel

	err := b.db.Get(&user, "SELECT * FROM users WHERE firebase_uid = $1;", UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (b *baseUserRepository) GetUserByUID(UID string) (*domain.UserModel, error) {
	var user domain.UserModel

	err := b.db.Get(&user, "SELECT * FROM users WHERE uid = $1;", UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (b *baseUserRepository) GetAdminByUserID(userID int) (*domain.AdminModel, error) {
	var admin domain.AdminModel

	err := b.db.Get(&admin, "SELECT * FROM admins WHERE user_id = $1;", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &admin, nil
}
