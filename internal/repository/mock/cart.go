package mock

import (
	"context"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/stretchr/testify/mock"
)

type CartRepositoryMock struct {
	mock.Mock
}

func (m *CartRepositoryMock) CreateCart(ctx context.Context) (domain.Cart, error) {
	args := m.Called(ctx)
	return args.Get(0).(domain.Cart), args.Error(1)
}

func (m *CartRepositoryMock) AddToCart(ctx context.Context, cartID int, product string, price float64) (domain.CartItem, error) {
	args := m.Called(ctx, cartID, product, price)
	return args.Get(0).(domain.CartItem), args.Error(1)
}

func (m *CartRepositoryMock) GetItemsByCartID(ctx context.Context, cartID int) ([]domain.CartItem, error) {
	args := m.Called(ctx, cartID)
	return args.Get(0).([]domain.CartItem), args.Error(1)

}
func (m *CartRepositoryMock) CartExists(ctx context.Context, cartID int) (bool, error) {
	args := m.Called(ctx, cartID)
	return args.Bool(0), args.Error(1)
}

func (m *CartRepositoryMock) DeleteCartItem(ctx context.Context, cartItemID int, cartID int) error {
	args := m.Called(ctx, cartItemID, cartID)
	return args.Error(0)
}
