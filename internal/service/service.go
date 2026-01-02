package service

import "github.com/SofiyaAndreyeva/cart-api/internal/repository"

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
