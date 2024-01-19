package usecase

import (
	"errors"
	"time"

	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
)

type baseAuthUsecase struct {
	authRepository domain.AuthRepository
	hashUtil       domain.HashUtil
	jwtUtil        domain.JWTUtil
}

func NewAuthUsecase(authRepository domain.AuthRepository, hashUtil domain.HashUtil, jwtUtil domain.JWTUtil) domain.AuthUsecase {
	return &baseAuthUsecase{authRepository: authRepository, hashUtil: hashUtil, jwtUtil: jwtUtil}
}

func (b *baseAuthUsecase) GetUserByUID(UID string) (*domain.UserModel, error) {
	user, err := b.authRepository.GetUserByUID(UID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b *baseAuthUsecase) GetAdminByUserID(userID int) (*domain.AdminModel, error) {
	admin, err := b.authRepository.GetAdminByUserID(userID)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (b *baseAuthUsecase) SignUp(email, password string) (accessToken, refreshToken string, accessTokenExpirationTime, refreshTokenExpirationTime time.Time, err error) {
	user, err := b.authRepository.GetUserByEmail(email)
	if user != nil {
		return "", "", time.Time{}, time.Time{}, errors.New("user already exist")
	}
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	hashedPassword, err := b.hashUtil.HashPassword(password)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}
	metadata := utils.GenerateMetadata()

	userPayload := &domain.AuthRepositoryPayloadCreateUser{
		UID:       metadata.UID(),
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: metadata.CreatedAt,
		UpdatedAt: metadata.UpdatedAt,
	}
	_, err = b.authRepository.CreateUser(userPayload)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	accessToken, accessTokenExpirationTime, err = b.jwtUtil.GenerateAccessToken(userPayload.UID)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}
	refreshToken, refreshTokenExpirationTime, err = b.jwtUtil.GenerateRefreshToken(userPayload.UID)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return
}

func (b *baseAuthUsecase) SignIn(email, password string) (accessToken, refreshToken string, accessTokenExpirationTime, refreshTokenExpirationTime time.Time, err error) {
	user, err := b.authRepository.GetUserByEmail(email)
	if user == nil {
		return "", "", time.Time{}, time.Time{}, errors.New("user not found")
	}
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	isPasswordCorrect := b.hashUtil.ValidatePassword(password, user.Password)
	if !isPasswordCorrect {
		return "", "", time.Time{}, time.Time{}, errors.New("wrong password")
	}

	accessToken, accessTokenExpirationTime, err = b.jwtUtil.GenerateAccessToken(user.UID)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}
	refreshToken, refreshTokenExpirationTime, err = b.jwtUtil.GenerateRefreshToken(user.UID)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return
}

func (b *baseAuthUsecase) RefreshAccessToken(refreshToken string) (string, time.Time, error) {
	newAccessToken, expirationTime, err := b.jwtUtil.Refresh(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	return newAccessToken, expirationTime, nil
}
