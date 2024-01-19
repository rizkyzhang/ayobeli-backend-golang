package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rizkyzhang/ayobeli-backend/domain"
)

type baseProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &baseProductRepository{db: db}
}

func (b *baseProductRepository) Create(productPayload *domain.ProductRepositoryPayloadCreateProduct) (string, error) {
	_, err := b.db.NamedExec(`
	INSERT INTO products (
    uid, name, slug, sku, description, images, weight, weight_value, base_price_value, base_price, offer_price_value, offer_price, discount, stock, status, created_at, updated_at
  )
	VALUES (
    :uid, :name, :slug, :sku, :description, :images, :weight, :weight_value,:base_price_value, :base_price, :offer_price_value, :offer_price, :discount, :stock, :status, :created_at, :updated_at
  );
	`, productPayload)
	if err != nil {
		return "", err
	}

	return productPayload.UID, nil
}

func (b *baseProductRepository) List(limit, cursor int, direction string) ([]*domain.ProductModel, error) {
	var products []*domain.ProductModel

	if direction == "" {
		err := b.db.Select(&products, `
			SELECT *
			FROM products
			LIMIT $1;
		`, limit)
		if err != nil {
			return nil, err
		}
	} else if direction == "next" {
		err := b.db.Select(&products, `
			SELECT *
			FROM products
			WHERE id > $1
			LIMIT $2;
		`, cursor, limit)
		if err != nil {
			return nil, nil
		}
	} else {
		err := b.db.Select(&products, `
			SELECT * 
			FROM (
				SELECT *
				FROM products
				WHERE id < $1
				ORDER by id DESC
				LIMIT $2
			) AS p
			ORDER by id ASC;
		`, cursor, limit)
		if err != nil {
			return nil, nil
		}
	}

	return products, nil
}

func (b *baseProductRepository) GetByUID(UID string) (*domain.ProductModel, error) {
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

func (b *baseProductRepository) UpdateByUID(productPayload *domain.ProductRepositoryPayloadUpdateProduct) error {
	_, err := b.db.NamedExec(`
  UPDATE products 
	SET name = :name,
			slug = :slug,
			sku = :sku,
			description = :description,
			images = :images,
			weight = :weight,
			weight_value = :weight_value,
      base_price_value = :base_price_value,
      base_price = :base_price,
      offer_price_value = :offer_price_value,
      offer_price = :offer_price,
			discount = :discount,
      stock = :stock,
      status = :status,
			updated_at = :updated_at
	WHERE uid = :uid;
	`, productPayload)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseProductRepository) DeleteByUID(UID string) error {
	_, err := b.db.Exec("DELETE FROM products WHERE uid = $1;", UID)
	if err != nil {
		return err
	}

	return nil
}
