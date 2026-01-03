package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
