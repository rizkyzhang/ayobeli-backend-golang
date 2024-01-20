package bootstrap

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
	"google.golang.org/api/option"
)

func NewFirebaseAuth(env *domain.Env) *auth.Client {
	opt := option.WithCredentialsFile(env.FirebaseCredentialPath)

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Failed to create Firebase app: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Firebase auth client: %v", err)
	}

	return authClient
}
