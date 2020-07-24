package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	userTableName = "users"
)

type Config struct {
	Host     string
	Port     string
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.SSLMode, cfg.Password,
	))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
