package service

import (
	"context"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/SofiyaAndreyeva/cart-api/internal/repository"
)

type Cart interface {
	CreateCart(ctx context.Context) (domain.Cart, error)
}

type Service struct {
	Cart
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Cart: NewCartService(repo),
	}
}
