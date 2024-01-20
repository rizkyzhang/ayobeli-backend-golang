package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

type baseJWTUtil struct {
	accessTokenSecretKey   []byte
	refreshTokenSecretKey  []byte
	accessTokenExpiryHour  int
	refreshTokenExpiryHour int
}

func NewJWTUtil(accessTokenSecretKey, refreshTokenSecretKey []byte, accessTokenExpiryHour, refreshTokenExpiryHour int) domain.JWTUtil {
	return &baseJWTUtil{accessTokenSecretKey: accessTokenSecretKey, refreshTokenSecretKey: refreshTokenSecretKey, accessTokenExpiryHour: accessTokenExpiryHour, refreshTokenExpiryHour: refreshTokenExpiryHour}
}

func (b *baseJWTUtil) generateToken(userUID string, tokenSecretKey []byte, tokenExpiryHour int) (string, time.Time, error) {
	expirationTime := time.Now().Add(time.Duration(tokenExpiryHour) * time.Hour)
	claims := &domain.JWTAccessTokenClaims{
		UserUID: userUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenSecretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func (b *baseJWTUtil) GenerateAccessToken(userUID string) (string, time.Time, error) {
	token, expirationTime, err := b.generateToken(userUID, b.accessTokenSecretKey, b.accessTokenExpiryHour)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expirationTime, err
}

func (b *baseJWTUtil) GenerateRefreshToken(userUID string) (string, time.Time, error) {
	token, expirationTime, err := b.generateToken(userUID, b.refreshTokenSecretKey, b.refreshTokenExpiryHour)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expirationTime, err
}

func (b *baseJWTUtil) ParseUserUID(tokenString string, isAccessToken bool) (string, error) {
	claims := &domain.JWTAccessTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if isAccessToken {
			return b.accessTokenSecretKey, nil
		}

		return b.refreshTokenSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", errors.New("invalid token signature")
		}

		return "", err
	}
	if !token.Valid {
		return "", err
	}

	if time.Until(claims.ExpiresAt.Time) < 30*time.Second {
		return "", errors.New("token already expired")
	}

	return claims.UserUID, nil
}

func (b *baseJWTUtil) Refresh(refreshToken string) (string, time.Time, error) {
	// Validate
	claims := &domain.JWTAccessTokenClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return b.refreshTokenSecretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", time.Time{}, errors.New("invalid token signature")
		}

		return "", time.Time{}, err
	}
	if !token.Valid {
		return "", time.Time{}, err
	}

	// Generate
	accessToken, expirationTime, err := b.generateToken(claims.UserUID, b.accessTokenSecretKey, b.accessTokenExpiryHour)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, expirationTime, nil
}
