package domain

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
)

type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &s)
}

// Controller
type ProductController interface {
	Create(c echo.Context) error
	GetByUID(c echo.Context) error
	UpdateByUID(c echo.Context) error
	DeleteByUID(c echo.Context) error
}

type ProductControllerPayloadCreateProduct struct {
	Name           string      `json:"name" validate:"required, min=5"`
	SKU            string      `json:"sku"`
	Description    string      `json:"description" validate:"required, min=30"`
	Images         StringSlice `json:"images" validate:"required, min=1"`
	WeightValue    float64     `json:"weight_value" validate:"required, min=100"`
	BasePriceValue int         `json:"base_price_value" validate:"required, min=5000"`
	Discount       *int        `json:"discount" validate:"required, max=100"`
	Stock          *int        `json:"stock" validate:"required"`
	Status         string      `json:"status" validate:"required, oneof=ACTIVE INACTIVE"`
}

type ProductControllerPayloadUpdateProduct struct {
	Name           string      `json:"name" validate:"required, min=5"`
	SKU            string      `json:"sku"`
	Description    string      `json:"description" validate:"required, min=30"`
	Images         StringSlice `json:"images" validate:"required, min=1"`
	WeightValue    float64     `json:"weight_value" validate:"required, min=100"`
	BasePriceValue int         `json:"base_price_value" validate:"required, min=5000"`
	Discount       *int        `json:"discount" validate:"required, max=100"`
	Stock          *int        `json:"stock" validate:"required"`
	Status         string      `json:"status" validate:"required, oneof=ACTIVE INACTIVE"`
}

type ProductControllerResponseGetProductByUID struct {
	UID             string      `json:"uid"`
	Name            string      `json:"name"`
	Slug            string      `json:"slug"`
	SKU             string      `json:"sku"`
	Description     string      `json:"description"`
	Images          StringSlice `json:"images"`
	Weight          string      `json:"weight"`
	WeightValue     float64     `json:"weight_value"`
	BasePrice       string      `json:"base_price"`
	BasePriceValue  int         `json:"base_price_value"`
	OfferPrice      string      `json:"offer_price"`
	OfferPriceValue int         `json:"offer_price_value"`
	Discount        int         `json:"discount"`
	Stock           int         `json:"stock"`
	Status          string      `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductControllerResponseListProducts struct {
	Products    []*ProductControllerResponseGetProductByUID
	IsFirstPage bool
	Limit       int
	PrevCursor  string
	NextCursor  string
}

// Usecase
type ProductUsecase interface {
	Create(payload *ProductUsecasePayloadCreateProduct) (string, error)
	List(limit int, encryptedCursor, direction string) (*ProductControllerResponseListProducts, error)
	GetByUID(UID string) (*ProductControllerResponseGetProductByUID, error)
	UpdateByUID(UID string, payload *ProductUsecasePayloadUpdateProduct) error
	DeleteByUID(UID string) error
}

type ProductUsecasePayloadCreateProduct struct {
	Name           string      `json:"name"`
	SKU            string      `json:"sku"`
	Description    string      `json:"description"`
	Images         StringSlice `json:"images"`
	WeightValue    float64     `json:"weight_value"`
	BasePriceValue int         `json:"base_price_value"`
	Discount       int         `json:"discount"`
	Stock          int         `json:"stock"`
	Status         string      `json:"status"`
}

type ProductUsecasePayloadUpdateProduct struct {
	Name           string      `json:"name"`
	SKU            string      `json:"sku"`
	Description    string      `json:"description"`
	Images         StringSlice `json:"images"`
	WeightValue    float64     `json:"weight_value"`
	BasePriceValue int         `json:"base_price_value"`
	Discount       int         `json:"discount"`
	Stock          int         `json:"stock"`
	Status         string      `json:"status"`
}

// Repository
type ProductModel struct {
	ID              int            `db:"id" json:"id"`
	UID             string         `db:"uid" json:"uid"`
	Name            string         `db:"name" json:"name"`
	Slug            string         `db:"slug" json:"slug"`
	SKU             sql.NullString `db:"sku" json:"sku"`
	Description     string         `db:"description" json:"description"`
	Images          StringSlice    `db:"images" json:"images"`
	Weight          string         `db:"weight" json:"weight"`
	WeightValue     float64        `db:"weight_value" json:"weight_value"`
	BasePrice       string         `db:"base_price" json:"base_price"`
	BasePriceValue  int            `db:"base_price_value" json:"base_price_value"`
	OfferPrice      string         `db:"offer_price" json:"offer_price"`
	OfferPriceValue int            `db:"offer_price_value" json:"offer_price_value"`
	Discount        int            `db:"discount" json:"discount"`
	Stock           int            `db:"stock" json:"stock"`
	Status          string         `db:"status" json:"status"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ProductRepository interface {
	Create(productPayload *ProductRepositoryPayloadCreateProduct) (string, error)
	List(limit, cursor int, direction string) ([]*ProductModel, error)
	GetByUID(UID string) (*ProductModel, error)
	UpdateByUID(productPayload *ProductRepositoryPayloadUpdateProduct) error
	DeleteByUID(UID string) error
}

type ProductRepositoryPayloadCreateProduct struct {
	UID             string      `db:"uid" json:"uid"`
	Name            string      `db:"name" json:"name"`
	Slug            string      `db:"slug" json:"slug"`
	SKU             string      `db:"sku" json:"sku"`
	Description     string      `db:"description" json:"description"`
	Images          StringSlice `db:"images" json:"images"`
	Weight          string      `db:"weight" json:"weight"`
	WeightValue     float64     `db:"weight_value" json:"weight_value"`
	BasePrice       string      `db:"base_price" json:"base_price"`
	BasePriceValue  int         `db:"base_price_value" json:"base_price_value"`
	OfferPrice      string      `db:"offer_price" json:"offer_price"`
	OfferPriceValue int         `db:"offer_price_value" json:"offer_price_value"`
	Discount        int         `db:"discount" json:"discount"`
	Stock           int         `db:"stock" json:"stock"`
	Status          string      `db:"status" json:"status"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ProductRepositoryPayloadUpdateProduct struct {
	UID             string      `db:"uid" json:"uid"`
	Name            string      `db:"name" json:"name"`
	Slug            string      `db:"slug" json:"slug"`
	SKU             string      `db:"sku" json:"sku"`
	Description     string      `db:"description" json:"description"`
	Images          StringSlice `db:"images" json:"images"`
	Weight          string      `db:"weight" json:"weight"`
	WeightValue     float64     `db:"weight_value" json:"weight_value"`
	BasePrice       string      `db:"base_price" json:"base_price"`
	BasePriceValue  int         `db:"base_price_value" json:"base_price_value"`
	OfferPrice      string      `db:"offer_price" json:"offer_price"`
	OfferPriceValue int         `db:"offer_price_value" json:"offer_price_value"`
	Discount        int         `db:"discount" json:"discount"`
	Stock           int         `db:"stock" json:"stock"`
	Status          string      `db:"status" json:"status"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
