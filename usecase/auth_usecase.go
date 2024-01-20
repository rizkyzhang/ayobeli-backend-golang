package usecase

import (
	"errors"

	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
)

type baseAuthUsecase struct {
	env            *domain.Env
	authRepository domain.AuthRepository
	authUtil       domain.AuthUtil
}

func NewAuthUsecase(env *domain.Env, authRepository domain.AuthRepository, authUtil domain.AuthUtil) domain.AuthUsecase {
	return &baseAuthUsecase{
		env:            env,
		authRepository: authRepository,
		authUtil:       authUtil,
	}
}

func (b *baseAuthUsecase) SignUp(email, password string) error {
	user, err := b.authRepository.GetUserByEmail(email)
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
	userPayload := &domain.AuthRepositoryPayloadCreateUser{
		UID:         metadata.UID(),
		FirebaseUID: firebaseUID,
		Email:       email,
		CreatedAt:   metadata.CreatedAt,
		UpdatedAt:   metadata.UpdatedAt,
	}
	_, err = b.authRepository.CreateUser(userPayload)
	if err != nil {
		return err
	}

	return nil
}

func (b *baseAuthUsecase) GetAccessToken(email, password string) (string, error) {
	user, err := b.authRepository.GetUserByEmail(email)
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

func (b *baseAuthUsecase) GetUserByFirebaseUID(UID string) (*domain.UserModel, error) {
	user, err := b.authRepository.GetUserByFirebaseUID(UID)
	if err != nil {
		return nil, err
	}

	return user, nil
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
