package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/rizkyzhang/ayobeli-backend/domain"
	"github.com/rizkyzhang/ayobeli-backend/internal/utils"
)

type Application struct {
	Env *domain.Env
	DB  *sqlx.DB
}

func App() Application {
	app := &Application{}
	app.Env = utils.LoadConfig("../.env")
	app.DB = NewPostgresDB(app.Env)
	return *app
}

func (app *Application) CloseDBConnection() {
	ClosePostgresDBConnection(app.DB)
}
