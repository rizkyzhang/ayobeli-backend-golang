package usecase

import (
	"context"

	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

type baseUserUsecase struct {
	env            *domain.Env
	userRepository domain.UserRepository
}

func NewUserUsecase(env *domain.Env, userRepository domain.UserRepository) domain.UserUsecase {
	return &baseUserUsecase{
		env:            env,
		userRepository: userRepository,
	}
}

func (b *baseUserUsecase) GetUserByFirebaseUID(ctx context.Context, UID string) (*domain.UserModel, error) {
	user, err := b.userRepository.GetUserByFirebaseUID(UID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b *baseUserUsecase) GetUserByUID(ctx context.Context, UID string) (*domain.UserModel, error) {
	user, err := b.userRepository.GetUserByUID(UID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (b *baseUserUsecase) GetAdminByUserID(ctx context.Context, userID int) (*domain.AdminModel, error) {
	admin, err := b.userRepository.GetAdminByUserID(userID)
	if err != nil {
		return nil, err
	}

	return admin, nil
}
