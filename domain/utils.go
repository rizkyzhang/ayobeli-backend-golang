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
