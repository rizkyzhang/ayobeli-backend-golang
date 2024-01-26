package usecase

import (
	"errors"

	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
)

type baseAuthUsecase struct {
	env            *domain.Env
	userRepository domain.UserRepository
	authUtil       domain.AuthUtil
}

func NewAuthUsecase(env *domain.Env, userRepository domain.UserRepository, authUtil domain.AuthUtil) domain.AuthUsecase {
	return &baseAuthUsecase{
		env:            env,
		userRepository: userRepository,
		authUtil:       authUtil,
	}
}

func (b *baseAuthUsecase) SignUp(email, password string, isAdmin bool) error {
	user, err := b.userRepository.GetUserByEmail(email)
	if user != nil {
		return errors.New("user already exist")
	}
	if err != nil {
		return err
	}

	firebaseUID, err := b.authUtil.CreateUser(email, password)
	if err != nil {
		return err
	}

	metadata := utils.GenerateMetadata()
	userPayload := &domain.UserRepositoryPayloadCreateUser{
		UID:         metadata.UID(),
		FirebaseUID: firebaseUID,
		Email:       email,
		IsAdmin:     isAdmin,
		CreatedAt:   metadata.CreatedAt,
		UpdatedAt:   metadata.UpdatedAt,
	}
	_, err = b.userRepository.CreateUser(userPayload)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseAuthUsecase) GetAccessToken(email, password string) (string, error) {
	user, err := b.userRepository.GetUserByEmail(email)
	if user == nil {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", err
	}

	accessToken, err := b.authUtil.GetAccessToken(email, password)
	if err != nil {
		return "", err
	}
	if accessToken == "" {
		return "", errors.New("invalid password")
	}

	return accessToken, nil
}
