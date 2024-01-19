package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rizkyzhang/ayobeli-backend/domain"
)

func SetupTestDB(env *domain.Env) (*dockertest.Pool, *dockertest.Resource, *sqlx.DB) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", env.TestDBUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", env.TestDBPassword),
			fmt.Sprintf("POSTGRES_DB=%s", env.DBName),
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf(env.TestDBUrl, hostAndPort)
	log.Println("Connecting to database on url: ", databaseUrl)

	// Tell docker to hard kill the container in 120 seconds
	resource.Expire(120)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	var db *sqlx.DB
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		_db, err := sqlx.Open("pgx", databaseUrl)
		if err != nil {
			return err
		}

		db = _db

		return _db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrations
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../../migrations")
	migrationDir := filepath.Join("file://" + basePath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationDir,
		env.DBName, driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	return pool, resource, db
}
