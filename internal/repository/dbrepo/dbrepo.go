package dbrepo

import (
	"database/sql"

	"github.com/adrialopezbou/bookings-go/internal/config"
	"github.com/adrialopezbou/bookings-go/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB: conn,
	}
}

type testDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewTestRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}