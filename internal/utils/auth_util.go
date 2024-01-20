package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

type baseAuthUtil struct {
	env          *domain.Env
	firebaseAuth *auth.Client
}

func NewAuthUtil(env *domain.Env, firebaseAuth *auth.Client) domain.AuthUtil {
	return &baseAuthUtil{env: env, firebaseAuth: firebaseAuth}
}

func (b *baseAuthUtil) CreateUser(email, password string) (authUID string, err error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	firebaseUserRecord, err := b.firebaseAuth.CreateUser(context.Background(), params)
	if err != nil {
		return "", err
	}

	return firebaseUserRecord.UID, nil
}

func (b *baseAuthUtil) VerifyToken(token string) (authUID string, err error) {
	parsedToken, err := b.firebaseAuth.VerifyIDTokenAndCheckRevoked(context.Background(), token)
	if err != nil {
		return "", err
	}

	return parsedToken.UID, nil
}

func (b *baseAuthUtil) GetAccessToken(email, password string) (accessToken string, err error) {
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
