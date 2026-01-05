package repository

import (
	"context"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Cart interface {
	CreateCart(ctx context.Context) (domain.Cart, error)
	AddToCart(ctx context.Context, cartID int, product string, price float64) (domain.CartItem, error)
	GetItemsByCartID(ctx context.Context, cartID int) ([]domain.CartItem, error)
	CartExists(ctx context.Context, cartID int) (bool, error)
}
type Repository struct {
	Cart
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Cart: NewCartRepository(db),
	}
}
