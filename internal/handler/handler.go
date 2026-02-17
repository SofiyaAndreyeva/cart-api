package handler

import (
	"context"
	"net/http"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
)


type CartServ interface {
	CreateCart(ctx context.Context) (domain.Cart, error)
	AddToCart(ctx context.Context, cartID int, product string, price float64) (domain.CartItem, error)
	DeleteFromCart(ctx context.Context, cartID int, cartItemID int) error
	GetCartItems(ctx context.Context, cartID int) (domain.Cart, error)
	GetCartPrice(ctx context.Context, cartID int) (domain.CartPriceResponse, error)
}

type Handler struct {
	service CartServ
}

func NewHandler(service CartServ) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitialRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /carts", h.createCart)                     
	mux.HandleFunc("POST /carts/{id}/items", h.addToCart)          
	mux.HandleFunc("DELETE /carts/{id}/items/{itemId}", h.deleteItemFromCart) 
	mux.HandleFunc("GET /carts/{id}", h.getCart)                  
	mux.HandleFunc("GET /carts/{id}/price", h.getCartPrice)        

	return mux
}
