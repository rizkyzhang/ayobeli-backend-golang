package bootstrap

import (
	"context"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rizkyzhang/ayobeli-backend-golang/domain"
)

func NewPostgresDB(env *domain.Env) *sqlx.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbUrl := env.DBUrl
	db, err := sqlx.ConnectContext(ctx, "pgx", dbUrl)
	if err != nil {
		log.Fatalf("Can't connect to Postgres DB with error %s", err)
	}

	return db
}

func ClosePostgresDBConnection(db *sqlx.DB) {
	db.Close()
	log.Println("Connection to Postgres DB closed")
}
