package repository

import (
	"context"
	"fmt"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) CreateCart(ctx context.Context) (domain.Cart, error) {
	var id int
	const query = `INSERT INTO carts DEFAULT VALUES RETURNING id`

	err := r.db.QueryRowContext(ctx, query).Scan(&id)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("repository: failed to create cart: %w", err)
	}
	return domain.Cart{
		ID:    id,
		Items: []domain.CartItem{},
	}, nil
}
