package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
	"github.com/SofiyaAndreyeva/cart-api/internal/repository/mock"
	"github.com/SofiyaAndreyeva/cart-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartService_CreateCart(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		mockSetup   func(m *mock.CartRepositoryMock)
		expected    domain.Cart
		expectErr   bool
		errContains string
	}{
		{
			name: "success",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CreateCart", ctx).
					Return(domain.Cart{ID: 1}, nil).
					Once()
			},
			expected: domain.Cart{ID: 1},
		},
		{
			name: "repo error",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CreateCart", ctx).
					Return(domain.Cart{}, errors.New("db error")).
					Once()
			},
			expectErr:   true,
			errContains: "service: create cart",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mock.CartRepositoryMock)
			if tt.mockSetup != nil {
				tt.mockSetup(repoMock)
			}

			cs := service.NewCartService(repoMock)
			cart, err := cs.CreateCart(ctx)

			if tt.expectErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Empty(t, cart)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, cart)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestCartService_AddToCart(t *testing.T) {
	ctx := context.Background()
	cartID := 1
	product := "Phone"
	price := 100.0

	tests := []struct {
		name        string
		product     string
		price       float64
		mockSetup   func(m *mock.CartRepositoryMock)
		expectErr   error
		errContains string
	}{
		{
			name:      "empty product",
			product:   "",
			price:     price,
			expectErr: domain.ErrEmptyProduct,
		},
		{
			name:      "invalid price",
			product:   product,
			price:     -10,
			expectErr: domain.ErrInvalidPrice,
		},
		{
			name:    "cart not found",
			product: product,
			price:   price,
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(false, nil)
			},
			expectErr: domain.ErrCartNotFound,
		},
		{
			name:    "success",
			product: product,
			price:   price,
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(true, nil)
				m.On("GetItemsByCartID", ctx, cartID).
					Return([]domain.CartItem{}, nil)
				m.On("AddToCart", ctx, cartID, product, price).
					Return(domain.CartItem{ID: 1}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mock.CartRepositoryMock)
			if tt.mockSetup != nil {
				tt.mockSetup(repoMock)
			}

			cs := service.NewCartService(repoMock)
			item, err := cs.AddToCart(ctx, cartID, tt.product, tt.price)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Empty(t, item)
			} else {
				require.NoError(t, err)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestCartService_DeleteFromCart(t *testing.T) {
	ctx := context.Background()
	cartID := 1
	itemID := 1

	tests := []struct {
		name        string
		mockSetup   func(m *mock.CartRepositoryMock)
		expectErr   error
		errContains string
	}{
		{
			name: "cart not found",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(false, nil)
			},
			expectErr: domain.ErrCartNotFound,
		},
		{
			name: "success",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(true, nil)
				m.On("DeleteCartItem", ctx, itemID, cartID).Return(nil)
			},
		},
		{
			name: "delete error",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(true, nil)
				m.On("DeleteCartItem", ctx, itemID, cartID).
					Return(errors.New("delete failed"))
			},
			expectErr: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mock.CartRepositoryMock)
			if tt.mockSetup != nil {
				tt.mockSetup(repoMock)
			}

			cs := service.NewCartService(repoMock)
			err := cs.DeleteFromCart(ctx, cartID, itemID)

			if tt.expectErr != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestCartService_GetCartItems(t *testing.T) {
	ctx := context.Background()
	cartID := 1

	tests := []struct {
		name        string
		mockSetup   func(m *mock.CartRepositoryMock)
		expectErr   error
		errContains string
	}{
		{
			name: "success",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(true, nil)
				m.On("GetItemsByCartID", ctx, cartID).
					Return([]domain.CartItem{}, nil)
			},
		},
		{
			name: "cart not found",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(false, nil)
			},
			expectErr: domain.ErrCartNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mock.CartRepositoryMock)
			if tt.mockSetup != nil {
				tt.mockSetup(repoMock)
			}

			cs := service.NewCartService(repoMock)
			cart, err := cs.GetCartItems(ctx, cartID)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Empty(t, cart)
			} else {
				require.NoError(t, err)
				assert.Equal(t, cartID, cart.ID)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestCartService_GetCartPrice(t *testing.T) {
	ctx := context.Background()
	cartID := 1

	tests := []struct {
		name      string
		mockSetup func(m *mock.CartRepositoryMock)
		expectErr error
	}{
		{
			name: "success with discount",
			mockSetup: func(m *mock.CartRepositoryMock) {
				items := []domain.CartItem{
					{Price: 2000},
					{Price: 1500},
					{Price: 1000},
					{Price: 800},
				}
				m.On("CartExists", ctx, cartID).Return(true, nil)
				m.On("GetItemsByCartID", ctx, cartID).Return(items, nil)
			},
		},
		{
			name: "cart not found",
			mockSetup: func(m *mock.CartRepositoryMock) {
				m.On("CartExists", ctx, cartID).Return(false, nil)
			},
			expectErr: domain.ErrCartNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mock.CartRepositoryMock)
			if tt.mockSetup != nil {
				tt.mockSetup(repoMock)
			}

			cs := service.NewCartService(repoMock)
			resp, err := cs.GetCartPrice(ctx, cartID)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Empty(t, resp)
			} else {
				require.NoError(t, err)
				assert.Equal(t, cartID, resp.CartID)
				assert.NotZero(t, resp.TotalPrice)
			}

			repoMock.AssertExpectations(t)
		})
	}
}
