package usecase

import (
	"context"
	"errors"

	"firebase.google.com/go/v4/auth"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
)

type baseAuthUsecase struct {
	authRepository domain.AuthRepository
	firebaseAuth   *auth.Client
}

func NewAuthUsecase(authRepository domain.AuthRepository, firebaseAuth *auth.Client) domain.AuthUsecase {
	return &baseAuthUsecase{authRepository: authRepository, firebaseAuth: firebaseAuth}
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

func (b *baseAuthUsecase) SignUp(email, password string) error {
	user, err := b.authRepository.GetUserByEmail(email)
	if user != nil {
		return errors.New("user already exist")
	}
	if err != nil {
		return err
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	firebaseUserRecord, err := b.firebaseAuth.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	metadata := utils.GenerateMetadata()
	userPayload := &domain.AuthRepositoryPayloadCreateUser{
		UID:         metadata.UID(),
		FirebaseUID: firebaseUserRecord.UID,
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
