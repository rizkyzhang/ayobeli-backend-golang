package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
)

type baseCartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) domain.CartRepository {
	return &baseCartRepository{db: db}
}

func (b *baseCartRepository) GetProductByUID(UID string) (*domain.ProductModel, error) {
	var product domain.ProductModel
	err := b.db.Get(&product, "SELECT * FROM products WHERE UID = $1;", UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &product, nil
}

func (b *baseCartRepository) CreateCart() error {
	metadata := utils.GenerateMetadata()

	cart := domain.CartModel{
		UID:              metadata.UID(),
		Quantity:         0,
		TotalPrice:       "Rp 0",
		TotalPriceValue:  0,
		TotalWeight:      "Rp 0",
		TotalWeightValue: 0,
		UserID:           1,
		CreatedAt:        metadata.CreatedAt,
		UpdatedAt:        metadata.UpdatedAt,
	}

	_, err := b.db.NamedExec(`
	INSERT INTO carts (uid, quantity, total_price, total_price_value, total_weight, total_weight_value, user_id, created_at, updated_at)
	VALUES (:uid, :quantity, :total_price, :total_price_value, :total_weight, :total_weight_value, :user_id, :created_at, :updated_at)
	`, cart)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCartRepository) GetCartByUID(UID string) (*domain.CartModel, error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		tx.Rollback()
	}()

	var cart domain.CartModel
	var cartItems []domain.CartItemModel

	err = tx.Get(&cart, "SELECT * FROM carts WHERE uid = $1", UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	err = tx.Select(&cartItems, "SELECT * FROM cart_items WHERE cart_id = $1", cart.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cart.CartItems = cartItems

	return &cart, nil
}

func (b *baseCartRepository) GetCartByUserID(userID int) (*domain.CartModel, error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		tx.Rollback()
	}()

	var cart domain.CartModel
	var cartItems []domain.CartItemModel

	err = tx.Get(&cart, "SELECT * FROM carts WHERE user_id = $1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	err = tx.Select(&cartItems, "SELECT * FROM cart_items WHERE cart_id = $1", cart.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	cart.CartItems = cartItems

	return &cart, nil
}

func (b *baseCartRepository) CreateCartItem(cartItemPayload domain.CartRepositoryPayloadCreateCartItem, cartPayload domain.CartRepositoryPayloadUpdateCart) (string, error) {
	tx, err := b.db.Beginx()
	if err != nil {
		return "", err
	}
	defer func() {
		tx.Rollback()
	}()

	_, err = tx.NamedExec(`
	INSERT INTO cart_items
	(uid, quantity, total_price, total_price_value, total_weight, total_weight_value, product_name, product_slug, product_image, product_weight, product_weight_value, base_price, base_price_value, offer_price, offer_price_value, discount, cart_id, product_id, created_at, updated_at)
	VALUES (:uid, :quantity, :total_price, :total_price_value, :total_weight, :total_weight_value, :product_name, :product_slug, :product_image, :product_weight, :product_weight_value, :base_price, :base_price_value, :offer_price, :offer_price_value, :discount, :cart_id, :product_id, :created_at, :updated_at);
	`, cartItemPayload)
	if err != nil {
		return "", err
	}

	_, err = tx.NamedExec(`
	UPDATE carts
	SET quantity = :quantity,
			total_price = :total_price,
			total_price_value = :total_price_value,
			total_weight = :total_weight,
			total_weight_value = :total_weight_value,
			updated_at = :updated_at
	WHERE uid = :uid;
	`, cartPayload)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return cartItemPayload.UID, nil
}

func (b *baseCartRepository) UpdateCartItem(cartItemPayload domain.CartRepositoryPayloadUpdateCartItem, cartPayload domain.CartRepositoryPayloadUpdateCart) error {
	tx, err := b.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	_, err = tx.NamedExec(`
	UPDATE cart_items 
	SET quantity = :quantity, 
			total_price = :total_price,
			total_price_value = :total_price_value,
			total_weight = :total_weight,
			total_weight_value = :total_weight_value,
			updated_at = :updated_at
	WHERE uid = :uid;
	`, cartItemPayload)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
	UPDATE carts 
	SET quantity = :quantity, 
			total_price = :total_price,
			total_price_value = :total_price_value,
			total_weight = :total_weight,
			total_weight_value = :total_weight_value,
			updated_at = :updated_at
	WHERE uid = :uid;
	`, cartPayload)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCartRepository) DeleteCartItemByUID(UID string, cartPayload domain.CartRepositoryPayloadUpdateCart) error {
	tx, err := b.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	_, err = tx.Exec("DELETE FROM cart_items WHERE uid = $1;", UID)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(`
	UPDATE carts 
	SET quantity = :quantity, 
			total_price = :total_price,
			total_price_value = :total_price_value,
			total_weight = :total_weight,
			total_weight_value = :total_weight_value,
			updated_at = :updated_at
	WHERE uid = :uid;
	`, cartPayload)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (b *baseCartRepository) GetCartItemByUID(UID string) (*domain.CartItemModel, error) {
	var cartItem domain.CartItemModel

	err := b.db.Get(&cartItem, "SELECT * FROM cart_items WHERE uid = $1", UID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &cartItem, nil
}

func (b *baseCartRepository) GetCartItemByProductID(productID int) (*domain.CartItemModel, error) {
	var cartItem domain.CartItemModel

	err := b.db.Get(&cartItem, "SELECT * FROM cart_items WHERE product_id = $1", productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &cartItem, nil
}
