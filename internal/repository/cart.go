package repository

import (
	"context"
	"fmt"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	err := r.db.QueryRowxContext(ctx, query).Scan(&id)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("repository: failed to create cart: %w", err)
	}
	return domain.Cart{
		ID:    id,
		Items: []domain.CartItem{},
	}, nil
}

func (r *CartRepository) AddToCart(ctx context.Context, cartID int, product string, price float64) (domain.CartItem, error) {
	var item domain.CartItem
	const query = `INSERT INTO cart_items (cart_id, product, price) VALUES ($1, $2, $3) RETURNING id, cart_id, product, price`

	err := r.db.QueryRowxContext(ctx, query, cartID, product, price).Scan(
		&item.ID,
		&item.CartID,
		&item.Product,
		&item.Price,
	)
	if err != nil {
		return domain.CartItem{}, fmt.Errorf("repository: failed to add item to cart %d: %w", cartID, err)
	}
	return item, nil
}

func (r *CartRepository) GetItemsByCartID(ctx context.Context, cartID int) ([]domain.CartItem, error) {
	var items []domain.CartItem
	const query = `SELECT id, cart_id, product, price FROM cart_items WHERE cart_id = $1`

	err := r.db.SelectContext(ctx, &items, query, cartID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get items by cart %d: %w", cartID, err)
	}
	
	return items, nil
}

func (r *CartRepository) CartExists(ctx context.Context, cartID int) (bool, error) {
	var exists bool
	const query = `SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)`

	err := r.db.GetContext(ctx, &exists, query, cartID)
	return exists, err
}

func (r *CartRepository) DeleteCartItem(ctx context.Context, cartItemID int, cartID int) error {
	const query = `DELETE FROM cart_items WHERE id = $1 AND cart_id = $2`

	res, err := r.db.ExecContext(ctx, query, cartItemID, cartID)
	if err != nil {
		return fmt.Errorf("repository: failed to delete cart item %d: %w", cartItemID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: rows affected error: %w", err)
	}

	if affected == 0 {
		return domain.ErrCartItemNotFound
	}

	return nil
}
