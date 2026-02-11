package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
)

func (h *Handler) carts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createCart(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
func (h *Handler) cartById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	switch r.Method {
	case http.MethodPost:
		h.addToCart(w, r)
	case http.MethodDelete:
		h.deleteItemFromCart(w, r)
	case http.MethodGet:
		if len(parts) == 2 {
			h.getCart(w, r)
			return
		}
		if len(parts) == 3 && parts[2] == "price" {
			h.getCartPrice(w, r)
			return
		}
		http.Error(w, "invalid path", http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// POST http://localhost:3000/carts
func (h *Handler) createCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cart, err := h.service.Cart.CreateCart(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(cart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// localhost:3000/carts/cartId/items
func (h *Handler) addToCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[2] != "items" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	cartID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid cart id", http.StatusBadRequest)
		return
	}
	var req domain.AddToCartRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := h.service.AddToCart(ctx, cartID, req.Product, req.Price)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DELETE http://localhost:3000/carts/1/items/1
func (h *Handler) deleteItemFromCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[2] != "items" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	cartID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid cart id", http.StatusBadRequest)
		return
	}
	cartItemID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "invalid cart item id", http.StatusBadRequest)
		return
	}
	err = h.service.DeleteFromCart(ctx, cartID, cartItemID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET http://localhost:3000/carts/1
func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	cartID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid cart id", http.StatusBadRequest)
		return
	}
	data, err := h.service.GetCartItems(ctx, cartID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// GET http://localhost:8080/carts/1/price
func (h *Handler) getCartPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	cartID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid cart id", http.StatusBadRequest)
		return
	}
	data, err := h.service.GetCartPrice(ctx, cartID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
