package bootstrap

import (
	"firebase.google.com/go/v4/auth"
	"github.com/jmoiron/sqlx"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
)

type Application struct {
	Env          *domain.Env
	DB           *sqlx.DB
	FirebaseAuth *auth.Client
}

func App() Application {
	app := &Application{}
	app.Env = utils.LoadConfig(".env")
	app.DB = NewPostgresDB(app.Env)
	app.FirebaseAuth = NewFirebaseAuth(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgresDBConnection(app.DB)
}
