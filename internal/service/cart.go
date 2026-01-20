package service

import (
	"context"
	"fmt"
	"math"

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

func (cs *CartService) GetCartItems(ctx context.Context, cartID int) (domain.Cart, error) {
	exists, err := cs.repo.CartExists(ctx, cartID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("service: cart check failed: %w", err)
	}
	if !exists {
		return domain.Cart{}, domain.ErrCartNotFound
	}
	items, err := cs.repo.GetItemsByCartID(ctx, cartID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("check cart: %w", err)
	}
	if items == nil {
		items = []domain.CartItem{}
	}
	return domain.Cart{
		ID:    cartID,
		Items: items,
	}, nil
}

func (cs *CartService) GetCartPrice(ctx context.Context, cartID int) (domain.CartPriceResponse, error) {
	exists, err := cs.repo.CartExists(ctx, cartID)
	if err != nil {
		return domain.CartPriceResponse{}, fmt.Errorf("service: cart check failed: %w", err)
	}
	if !exists {
		return domain.CartPriceResponse{}, domain.ErrCartNotFound
	}

	items, err := cs.repo.GetItemsByCartID(ctx, cartID)
	if err != nil {
		return domain.CartPriceResponse{}, fmt.Errorf("check cart: %w", err)
	}

	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.Price
	}

	discount := 0

	if len(items) > 3 {
		discount = 5
	}

	if totalPrice > 5000 {
		discount = 10
	}

	finalPrice := totalPrice
	if discount > 0 {
		finalPrice = totalPrice - ((totalPrice * float64(discount)) / 100)
	}
	totalPrice = math.Round(totalPrice*100) / 100
	finalPrice = math.Round(finalPrice*100) / 100

	return domain.CartPriceResponse{
		CartID:          cartID,
		TotalPrice:      totalPrice,
		DiscountPercent: discount,
		FinalPrice:      finalPrice,
	}, nil
}
