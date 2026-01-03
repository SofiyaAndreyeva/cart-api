package repository

import (
	"fmt"

	"github.com/SofiyaAndreyeva/cart-api/internal/config"
	"github.com/jmoiron/sqlx"
)

func NewPostgresDB(cfg config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPass, cfg.DBSSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
