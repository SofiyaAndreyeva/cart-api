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
