package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout         int    `mapstructure:"CONTEXT_TIMEOUT"`
	TestDBUrl              string `mapstructure:"TEST_DB_URL"`
	TestDBUser             string `mapstructure:"TEST_DB_USER"`
	TestDBPassword         string `mapstructure:"TEST_DB_PASSWORD"`
	DBUrl                  string `mapstructure:"DB_URL"`
	DBName                 string `mapstructure:"DB_NAME"`
	AesSecret              string `mapstructure:"AES_SECRET"`
	AccessTokenExpiryHour  int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	RefreshTokenExpiryHour int    `mapstructure:"REFRESH_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret      string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret     string `mapstructure:"REFRESH_TOKEN_SECRET"`
}

type AesEncryptUtil interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type HashUtil interface {
	HashPassword(password string) (string, error)
	ValidatePassword(password, hash string) bool
}

type JWTAccessTokenClaims struct {
	UserUID string `json:"user_uid"`
	jwt.RegisteredClaims
}

type JWTUtil interface {
	GenerateAccessToken(userUID string) (string, time.Time, error)
	GenerateRefreshToken(userUID string) (string, time.Time, error)
	ParseUserUID(tokenString string, isAccessToken bool) (string, error)
	Refresh(refreshToken string) (string, time.Time, error)
}

type Metadata struct {
	UID       func() string
	Slug      func(str string) string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CalculatedPrice struct {
	Base       string
	Offer      string
	OfferValue int
}

type ProductUtil interface {
	CalculatePrice(baseValue int, discount int) (*CalculatedPrice, error)
	FormatRupiah(value int) (string, error)
	FormatWeight(weightInGram float64) string
}

type CalculatedCart struct {
	CartQuantity             int
	CartTotalPriceValue      int
	CartTotalPrice           string
	CartTotalWeightValue     float64
	CartTotalWeight          string
	CartItemTotalPriceValue  int
	CartItemTotalPrice       string
	CartItemTotalWeightValue float64
	CartItemTotalWeight      string
}

type CartUtil interface {
	CalculateCreateCartItem(payload *CartUsecasePayloadCreateCartItem) (*CalculatedCart, error)
	CalculateUpdateCartItem(payload *CartUsecasePayloadUpdateCartItem) (*CalculatedCart, error)
	CalculateDeleteCartItem(payload *CartUsecasePayloadDeleteCartItem) (*CalculatedCart, error)
}
