package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Config struct {
	DSN    string `json:"dsn"`
	Driver string `json:"driver"`
}

func NewConn(cfg *Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("db: nil config")
	}
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("db: sql open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db: ping: %w", err)
	}
	return db, nil
}
