package repository

import (
	"context"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Cart interface {
	CreateCart(ctx context.Context) (domain.Cart, error)
}
type Repository struct {
	Cart
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Cart: NewCartRepository(db),
	}
}
