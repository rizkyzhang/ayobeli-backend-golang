package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"github.com/rizkyzhang/ayobeli-backend-golang/internal/utils"
)

type baseAuthUsecase struct {
	env            *domain.Env
	authRepository domain.AuthRepository
	firebaseAuth   *auth.Client
}

func NewAuthUsecase(env *domain.Env, authRepository domain.AuthRepository, firebaseAuth *auth.Client) domain.AuthUsecase {
	return &baseAuthUsecase{
		env:            env,
		authRepository: authRepository,
		firebaseAuth:   firebaseAuth,
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

func (b *baseAuthUsecase) GetAccessToken(email, password string) (string, error) {
	user, err := b.authRepository.GetUserByEmail(email)
	if user == nil {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", err
	}

	reqBody := map[string]string{
		"email":             email,
		"password":          password,
		"returnSecureToken": "true",
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", b.env.FirebaseVerifyPasswordURL, bytes.NewBuffer(reqBytes))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	_resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}
	var resBody struct {
		IdToken string `json:"idToken"`
	}
	err = json.Unmarshal(_resBody, &resBody)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return "", err
	}

	return resBody.IdToken, nil
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
