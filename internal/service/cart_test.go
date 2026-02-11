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

func TestCartService_CreateCart_Success(t *testing.T) {
	ctx := context.Background()

	expectedCart := domain.Cart{
		ID:    1,
		Items: []domain.CartItem{},
	}

	repoMock := new(mock.CartRepositoryMock)

	repoMock.
		On("CreateCart", ctx).
		Return(expectedCart, nil).
		Once()

	cs := service.NewCartService(repoMock)

	cart, err := cs.CreateCart(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedCart, cart)

	repoMock.AssertExpectations(t)
}

func TestCartService_CreateCart_RepoError(t *testing.T) {
	ctx := context.Background()

	repoErr := errors.New("db error")

	repoMock := new(mock.CartRepositoryMock)
	repoMock.
		On("CreateCart", ctx).
		Return(domain.Cart{}, repoErr).
		Once()

	cs := service.NewCartService(repoMock)

	cart, err := cs.CreateCart(ctx)

	require.Error(t, err)
	assert.Empty(t, cart)

	assert.ErrorContains(t, err, "service: create cart")
	assert.ErrorIs(t, err, repoErr)

	repoMock.AssertExpectations(t)
}

func TestCartService_AddItems_Success(t *testing.T) {
	ctx := context.Background()
	cartID := 1
	product := "Phone"
	price := 121.32

	expectedItem := domain.CartItem{
		ID:      1,
		CartID:  cartID,
		Product: product,
		Price:   price,
	}

	repoMock := new(mock.CartRepositoryMock)
	repoMock.On("CartExists", ctx, cartID).Return(true, nil)
	repoMock.On("GetItemsByCartID", ctx, cartID).Return([]domain.CartItem{}, nil)
	repoMock.On("AddToCart", ctx, cartID, product, price).Return(expectedItem, nil)
	cs := service.NewCartService(repoMock)
	addedItem, err := cs.AddToCart(ctx, cartID, product, price)

	require.NoError(t, err)
	assert.Equal(t, expectedItem, addedItem)

	repoMock.AssertExpectations(t)
}

func TestCartService_AddToCart_Errors(t *testing.T) {
	ctx := context.Background()
	cartID := 1
	product := "Phone"
	price := 121.32

	t.Run("empty product", func(t *testing.T) {
		cs := service.NewCartService(new(mock.CartRepositoryMock))
		item, err := cs.AddToCart(ctx, cartID, "", price)

		require.ErrorIs(t, err, domain.ErrEmptyProduct)
		assert.Empty(t, item)
	})

	t.Run("invalid price", func(t *testing.T) {
		cs := service.NewCartService(new(mock.CartRepositoryMock))
		item, err := cs.AddToCart(ctx, cartID, product, -5)

		require.ErrorIs(t, err, domain.ErrInvalidPrice)
		assert.Empty(t, item)
	})

	t.Run("cart not found", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)
		repoMock.On("CartExists", ctx, cartID).Return(false, nil)

		cs := service.NewCartService(repoMock)
		item, err := cs.AddToCart(ctx, cartID, product, price)

		require.ErrorIs(t, err, domain.ErrCartNotFound)
		assert.Empty(t, item)
		repoMock.AssertExpectations(t)
	})

	t.Run("cart exists error", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)
		repoMock.On("CartExists", ctx, cartID).Return(false, errors.New("db error"))

		cs := service.NewCartService(repoMock)
		item, err := cs.AddToCart(ctx, cartID, product, price)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "service: cart check failed")
		assert.Empty(t, item)
		repoMock.AssertExpectations(t)
	})

	t.Run("cart limit exceeded", func(t *testing.T) {
		existingItems := []domain.CartItem{{}, {}, {}, {}, {}}

		repoMock := new(mock.CartRepositoryMock)
		repoMock.On("CartExists", ctx, cartID).Return(true, nil)
		repoMock.On("GetItemsByCartID", ctx, cartID).Return(existingItems, nil)

		cs := service.NewCartService(repoMock)
		item, err := cs.AddToCart(ctx, cartID, product, price)

		require.ErrorIs(t, err, domain.ErrLimitCart)
		assert.Empty(t, item)
		repoMock.AssertExpectations(t)
	})

	t.Run("add to db error", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)
		repoMock.On("CartExists", ctx, cartID).Return(true, nil)
		repoMock.On("GetItemsByCartID", ctx, cartID).Return([]domain.CartItem{}, nil)
		repoMock.On("AddToCart", ctx, cartID, product, price).Return(domain.CartItem{}, errors.New("db error"))

		cs := service.NewCartService(repoMock)
		item, err := cs.AddToCart(ctx, cartID, product, price)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "add to db")
		assert.Empty(t, item)
		repoMock.AssertExpectations(t)
	})
}

func TestCartService_DeleteCartItem_Success(t *testing.T) {
	ctx := context.Background()
	cartID := 1
	cartItemID := 1
	repoMock := new(mock.CartRepositoryMock)
	repoMock.On("CartExists", ctx, cartID).Return(true, nil)
	repoMock.On("DeleteCartItem", ctx, cartItemID, cartID).Return(nil)

	cs := service.NewCartService(repoMock)
	err := cs.DeleteFromCart(ctx, cartID, cartItemID)

	require.NoError(t, err)

	repoMock.AssertExpectations(t)
}

func TestCartService_DeleteCartItem_Error(t *testing.T) {
	ctx := context.Background()
	cartID := 1
	cartItemID := 1
	t.Run("cart not found", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)
		repoMock.On("CartExists", ctx, cartID).Return(false, nil)

		cs := service.NewCartService(repoMock)
		err := cs.DeleteFromCart(ctx, cartID, cartItemID)

		require.ErrorIs(t, err, domain.ErrCartNotFound)
		repoMock.AssertExpectations(t)
	})

	t.Run("cart exists error", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)
		repoMock.On("CartExists", ctx, cartID).Return(false, errors.New("db error"))

		cs := service.NewCartService(repoMock)
		err := cs.DeleteFromCart(ctx, cartID, cartItemID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "service: cart check failed")
		repoMock.AssertExpectations(t)
	})

	t.Run("cart item not found", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)

		repoMock.On("CartExists", ctx, cartID).Return(true, nil)
		repoMock.On("DeleteCartItem", ctx, cartItemID, cartID).
			Return(domain.ErrCartItemNotFound)

		cs := service.NewCartService(repoMock)
		err := cs.DeleteFromCart(ctx, cartID, cartItemID)

		require.ErrorIs(t, err, domain.ErrCartItemNotFound)
		repoMock.AssertExpectations(t)
	})
	t.Run("delete item error", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)
		repoMock.On("CartExists", ctx, cartID).Return(true, nil)
		repoMock.On("DeleteCartItem", ctx, cartItemID, cartID).
			Return(errors.New("delete failed"))

		cs := service.NewCartService(repoMock)
		err := cs.DeleteFromCart(ctx, cartID, cartItemID)

		require.Error(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestCartService_GetCartItems_Success(t *testing.T) {
	ctx := context.Background()
	cartID := 1

	repoMock := new(mock.CartRepositoryMock)

	repoMock.
		On("CartExists", ctx, cartID).
		Return(true, nil)

	repoMock.
		On("GetItemsByCartID", ctx, cartID).
		Return([]domain.CartItem{}, nil)

	cs := service.NewCartService(repoMock)

	cart, err := cs.GetCartItems(ctx, cartID)

	require.NoError(t, err)
	require.Equal(t, cartID, cart.ID)
	require.NotNil(t, cart.Items)
	require.Len(t, cart.Items, 0)

	repoMock.AssertExpectations(t)
}

func TestCartService_GetCartItems_Error(t *testing.T) {
	ctx := context.Background()
	cartID := 1

	t.Run("error from CartExists", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)

		repoMock.
			On("CartExists", ctx, cartID).
			Return(false, errors.New("db error"))

		cs := service.NewCartService(repoMock)

		cart, err := cs.GetCartItems(ctx, cartID)

		require.Error(t, err)
		require.Empty(t, cart)
		require.Contains(t, err.Error(), "service: cart check failed")

		repoMock.AssertExpectations(t)
	})

	t.Run("cart not found", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)

		repoMock.
			On("CartExists", ctx, cartID).
			Return(false, nil)

		cs := service.NewCartService(repoMock)

		cart, err := cs.GetCartItems(ctx, cartID)

		require.Error(t, err)
		require.Empty(t, cart)
		require.ErrorIs(t, err, domain.ErrCartNotFound)

		repoMock.AssertExpectations(t)
	})

	t.Run("error from GetItemsByCartID", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)

		repoMock.
			On("CartExists", ctx, cartID).
			Return(true, nil)

		repoMock.
			On("GetItemsByCartID", ctx, cartID).
			Return([]domain.CartItem(nil), errors.New("items error"))

		cs := service.NewCartService(repoMock)

		cart, err := cs.GetCartItems(ctx, cartID)

		require.Error(t, err)
		require.Empty(t, cart)
		require.Contains(t, err.Error(), "check cart")

		repoMock.AssertExpectations(t)
	})
}

func TestCartService_GetCartPrice_Success(t *testing.T) {
	ctx := context.Background()
	cartID := 1

	items := []domain.CartItem{
		{Price: 2000},
		{Price: 1500},
		{Price: 1000},
		{Price: 800},
	}

	repoMock := new(mock.CartRepositoryMock)

	repoMock.
		On("CartExists", ctx, cartID).
		Return(true, nil)

	repoMock.
		On("GetItemsByCartID", ctx, cartID).
		Return(items, nil)

	cs := service.NewCartService(repoMock)

	resp, err := cs.GetCartPrice(ctx, cartID)

	require.NoError(t, err)
	require.Equal(t, cartID, resp.CartID)
	require.Equal(t, 5300.0, resp.TotalPrice)
	require.Equal(t, 10, resp.DiscountPercent)
	require.Equal(t, 4770.0, resp.FinalPrice)

	repoMock.AssertExpectations(t)
}
func TestCartService_GetCartPrice_Error(t *testing.T) {
	ctx := context.Background()
	cartID := 1

	t.Run("error from CartExists", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)

		repoMock.
			On("CartExists", ctx, cartID).
			Return(false, errors.New("db error"))

		cs := service.NewCartService(repoMock)

		resp, err := cs.GetCartPrice(ctx, cartID)

		require.Error(t, err)
		require.Empty(t, resp)
		require.Contains(t, err.Error(), "service: cart check failed")

		repoMock.AssertExpectations(t)
	})

	t.Run("cart not found", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)

		repoMock.
			On("CartExists", ctx, cartID).
			Return(false, nil)

		cs := service.NewCartService(repoMock)

		resp, err := cs.GetCartPrice(ctx, cartID)

		require.Error(t, err)
		require.Empty(t, resp)
		require.ErrorIs(t, err, domain.ErrCartNotFound)

		repoMock.AssertExpectations(t)
	})

	t.Run("error from GetItemsByCartID", func(t *testing.T) {
		repoMock := new(mock.CartRepositoryMock)

		repoMock.
			On("CartExists", ctx, cartID).
			Return(true, nil)

		repoMock.
			On("GetItemsByCartID", ctx, cartID).
			Return([]domain.CartItem(nil), errors.New("items error"))

		cs := service.NewCartService(repoMock)

		resp, err := cs.GetCartPrice(ctx, cartID)

		require.Error(t, err)
		require.Empty(t, resp)
		require.Contains(t, err.Error(), "check cart")

		repoMock.AssertExpectations(t)
	})
}
