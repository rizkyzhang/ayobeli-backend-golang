package utils

import (
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"golang.org/x/crypto/bcrypt"
)

type baseHashUtil struct {
}

func NewHashUtil() domain.HashUtil {
	return &baseHashUtil{}
}

func (b *baseHashUtil) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (b *baseHashUtil) ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
