package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
)

type baseAuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) domain.AuthRepository {
	return &baseAuthRepository{db: db}
}

func (b *baseAuthRepository) CreateUser(userPayload *domain.AuthRepositoryPayloadCreateUser) (uint64, error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer func() {
		tx.Rollback()
	}()

	query, args, err := tx.BindNamed(`
	INSERT INTO users (uid, email, password, name, phone, profile_image, created_at, updated_at)
	VALUES (:uid, :email, :password, :name, :phone, :profile_image, :created_at, :updated_at)
	RETURNING id;
	`, userPayload)
	if err != nil {
		return 0, err
	}

	var userID uint64
	err = tx.Get(&userID, query, args...)
	if err != nil {
		return 0, err
	}

	metadata := utils.GenerateMetadata()

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

func (b *baseAuthRepository) GetUserByEmail(email string) (*domain.UserModel, error) {
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

func (b *baseAuthRepository) GetUserByUID(UID string) (*domain.UserModel, error) {
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

func (b *baseAuthRepository) GetAdminByUserID(userID uint64) (*domain.AdminModel, error) {
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
