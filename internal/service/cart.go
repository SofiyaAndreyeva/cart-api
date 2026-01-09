package service

import (
	"context"
	"fmt"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/SofiyaAndreyeva/cart-api/internal/repository"
)

type CartService struct {
	repo repository.Cart
}

func NewCartService(repo repository.Cart) *CartService {
	return &CartService{repo: repo}
}
func (cs *CartService) CreateCart(ctx context.Context) (domain.Cart, error) {
	cart, err := cs.repo.CreateCart(ctx)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("service: create cart: %w", err)
	}
	return cart, nil
}
func (cs *CartService) AddToCart(ctx context.Context, cartID int, product string, price float64) (domain.CartItem, error) {
	if product == "" {
		return domain.CartItem{}, domain.ErrEmptyProduct
	}
	if price <= 0 {
		return domain.CartItem{}, domain.ErrInvalidPrice
	}
	exists, err := cs.repo.CartExists(ctx, cartID)
	if err != nil {
		return domain.CartItem{}, fmt.Errorf("service: cart check failed: %w", err)
	}
	if !exists {
		return domain.CartItem{}, domain.ErrCartNotFound
	}
	items, err := cs.repo.GetItemsByCartID(ctx, cartID)
	if err != nil {
		return domain.CartItem{}, fmt.Errorf("check cart: %w", err)
	}
	if len(items) >= 5 {
		return domain.CartItem{}, domain.ErrLimitCart
	}
	addedItem, err := cs.repo.AddToCart(ctx, cartID, product, price)
	if err != nil {
		return domain.CartItem{}, fmt.Errorf("add to db: %w", err)
	}

	return addedItem, nil
}

func (cs *CartService) DeleteFromCart(ctx context.Context, cartID int, cartItemID int) error {
	existsCart, err := cs.repo.CartExists(ctx, cartID)
	if err != nil {
		return fmt.Errorf("service: cart check failed: %w", err)
	}
	if !existsCart {
		return domain.ErrCartNotFound
	}
	return cs.repo.DeleteCartItem(ctx, cartItemID, cartID)
}
