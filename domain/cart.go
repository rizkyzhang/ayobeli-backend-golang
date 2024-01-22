package domain

import (
	"time"

	"github.com/labstack/echo/v4"
)

// Controller
type CartController interface {
	// Cart
	GetCartByUserID(c echo.Context) error

	// Cart item
	CreateCartItem(c echo.Context) error
	UpdateCartItem(c echo.Context) error
	DeleteCartItemByUID(c echo.Context) error
}

type CartControllerPayloadCreateCartItem struct {
	ProductUID string `json:"product_uid"`
	Quantity   int    `json:"quantity"`
}

type CartControllerPayloadUpdateCartItem struct {
	Quantity int `json:"quantity"`
}

type CartControllerResponseGetCart struct {
	UID              string                               `db:"uid" json:"uid"`
	Quantity         int                                  `db:"quantity" json:"quantity"`
	TotalPrice       string                               `db:"total_price" json:"total_price"`
	TotalPriceValue  int                                  `db:"total_price_value" json:"total_price_value"`
	TotalWeight      string                               `db:"total_weight" json:"total_weight"`
	TotalWeightValue float64                              `db:"total_weight_value" json:"total_weight_value"`
	CartItems        []ControllerResponsePropertyCartItem `db:"cart_items" json:"cart_items"`
}

type ControllerResponsePropertyCartItem struct {
	UID              string  `db:"uid" json:"uid"`
	Quantity         int     `db:"quantity" json:"quantity"`
	TotalPrice       string  `db:"total_price" json:"total_price"`
	TotalPriceValue  int     `db:"total_price_value" json:"total_price_value"`
	TotalWeight      string  `db:"total_weight" json:"total_weight"`
	TotalWeightValue float64 `db:"total_weight_value" json:"total_weight_value"`

	// Product information
	ProductName        string  `db:"product_name" json:"product_name"`
	ProductSlug        string  `db:"product_slug" json:"product_slug"`
	ProductImage       string  `db:"product_image" json:"product_image"`
	ProductWeight      string  `db:"product_weight" json:"product_weight"`
	ProductWeightValue float64 `db:"product_weight_value" json:"product_weight_value"`
	BasePrice          string  `db:"base_price" json:"base_price"`
	BasePriceValue     int     `db:"base_price_value" json:"base_price_value"`
	OfferPrice         string  `db:"offer_price" json:"offer_price"`
	OfferPriceValue    int     `db:"offer_price_value" json:"offer_price_value"`
	Discount           uint8   `db:"discount" json:"discount"`
}

// Usecase
type CartUsecase interface {
	// Cart
	GetCartByUserID(userID int) (*CartControllerResponseGetCart, error)
	GetCartByUserIDMiddleware(userID int) (*CartModel, error)

	// Cart item
	CreateCartItem(payload *CartUsecasePayloadCreateCartItem) (string, error)
	GetCartItemByUID(UID string) (*CartItemModel, error)
	GetCartItemByProductID(productID int) (*CartItemModel, error)
	UpdateCartItem(payload *CartUsecasePayloadUpdateCartItem) error
	DeleteCartItemByUID(payload *CartUsecasePayloadDeleteCartItem) error
}

type CartUsecasePayloadCreateCartItem struct {
	Cart     *CartModel    `json:"cart"`
	Product  *ProductModel `json:"product"`
	Quantity int           `json:"quantity"`
}

type CartUsecasePayloadUpdateCartItem struct {
	Cart     *CartModel     `json:"cart"`
	CartItem *CartItemModel `json:"cart_item"`
	UID      string         `json:"uid"`
	Quantity int            `json:"quantity"`
}

type CartUsecasePayloadDeleteCartItem struct {
	Cart     *CartModel     `json:"cart"`
	CartItem *CartItemModel `json:"cart_item"`
	UID      string         `json:"uid"`
}

// Repository
type CartModel struct {
	ID               int     `db:"id" json:"id"`
	UID              string  `db:"uid" json:"uid"`
	Quantity         int     `db:"quantity" json:"quantity"`
	TotalPrice       string  `db:"total_price" json:"total_price"`
	TotalPriceValue  int     `db:"total_price_value" json:"total_price_value"`
	TotalWeight      string  `db:"total_weight" json:"total_weight"`
	TotalWeightValue float64 `db:"total_weight_value" json:"total_weight_value"`

	// Relationship
	CartItems []CartItemModel `db:"cart_items" json:"cart_items"`
	UserID    int             `db:"user_id" json:"user_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CartItemModel struct {
	ID               int     `db:"id" json:"id"`
	UID              string  `db:"uid" json:"uid"`
	Quantity         int     `db:"quantity" json:"quantity"`
	TotalPrice       string  `db:"total_price" json:"total_price"`
	TotalPriceValue  int     `db:"total_price_value" json:"total_price_value"`
	TotalWeight      string  `db:"total_weight" json:"total_weight"`
	TotalWeightValue float64 `db:"total_weight_value" json:"total_weight_value"`

	// Product information
	ProductName        string  `db:"product_name" json:"product_name"`
	ProductSlug        string  `db:"product_slug" json:"product_slug"`
	ProductImage       string  `db:"product_image" json:"product_image"`
	ProductWeight      string  `db:"product_weight" json:"product_weight"`
	ProductWeightValue float64 `db:"product_weight_value" json:"product_weight_value"`
	BasePrice          string  `db:"base_price" json:"base_price"`
	BasePriceValue     int     `db:"base_price_value" json:"base_price_value"`
	OfferPrice         string  `db:"offer_price" json:"offer_price"`
	OfferPriceValue    int     `db:"offer_price_value" json:"offer_price_value"`
	Discount           int     `db:"discount" json:"discount"`

	// Relationship
	CartID    int `db:"cart_id" json:"cart_id"`
	ProductID int `db:"product_id" json:"product_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CartRepository interface {
	GetProductByUID(UID string) (*ProductModel, error)

	// Cart
	CreateCart() error
	GetCartByUID(UID string) (*CartModel, error)
	GetCartByUserID(userID int) (*CartModel, error)

	// Cart item
	CreateCartItem(cartItemPayload CartRepositoryPayloadCreateCartItem, cartPayload CartRepositoryPayloadUpdateCart) (string, error)
	GetCartItemByUID(UID string) (*CartItemModel, error)
	GetCartItemByProductID(productID int) (*CartItemModel, error)
	UpdateCartItem(cartItemPayload CartRepositoryPayloadUpdateCartItem, cartPayload CartRepositoryPayloadUpdateCart) error
	DeleteCartItemByUID(UID string, cartPayload CartRepositoryPayloadUpdateCart) error
}

type CartRepositoryPayloadUpdateCart struct {
	UID              string  `db:"uid" json:"uid"`
	Quantity         int     `db:"quantity" json:"quantity"`
	TotalPrice       string  `db:"total_price" json:"total_price"`
	TotalPriceValue  int     `db:"total_price_value" json:"total_price_value"`
	TotalWeight      string  `db:"total_weight" json:"total_weight"`
	TotalWeightValue float64 `db:"total_weight_value" json:"total_weight_value"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CartRepositoryPayloadCreateCartItem struct {
	UID              string  `db:"uid" json:"uid"`
	Quantity         int     `db:"quantity" json:"quantity"`
	TotalPrice       string  `db:"total_price" json:"total_price"`
	TotalPriceValue  int     `db:"total_price_value" json:"total_price_value"`
	TotalWeight      string  `db:"total_weight" json:"total_weight"`
	TotalWeightValue float64 `db:"total_weight_value" json:"total_weight_value"`

	// Product information
	ProductName        string  `db:"product_name" json:"product_name"`
	ProductSlug        string  `db:"product_slug" json:"product_slug"`
	ProductImage       string  `db:"product_image" json:"product_image"`
	ProductWeight      string  `db:"product_weight" json:"product_weight"`
	ProductWeightValue float64 `db:"product_weight_value" json:"product_weight_value"`
	BasePrice          string  `db:"base_price" json:"base_price"`
	BasePriceValue     int     `db:"base_price_value" json:"base_price_value"`
	OfferPrice         string  `db:"offer_price" json:"offer_price"`
	OfferPriceValue    int     `db:"offer_price_value" json:"offer_price_value"`
	Discount           int     `db:"discount" json:"discount"`

	// Relationship
	CartID    int `db:"cart_id" json:"cart_id"`
	ProductID int `db:"product_id" json:"product_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CartRepositoryPayloadUpdateCartItem struct {
	UID              string  `db:"uid" json:"uid"`
	Quantity         int     `db:"quantity" json:"quantity"`
	TotalPrice       string  `db:"total_price" json:"total_price"`
	TotalPriceValue  int     `db:"total_price_value" json:"total_price_value"`
	TotalWeight      string  `db:"total_weight" json:"total_weight"`
	TotalWeightValue float64 `db:"total_weight_value" json:"total_weight_value"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
