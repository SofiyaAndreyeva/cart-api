package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
)


// POST http://localhost:3000/carts
func (h *Handler) createCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cart, err := h.service.CreateCart(ctx)
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


// POST localhost:3000/carts/cartId/items
func (h *Handler) addToCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cartID, err := strconv.Atoi(r.PathValue("id"))
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

	cartID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid cart id", http.StatusBadRequest)
		return
	}

	cartItemID, err := strconv.Atoi(r.PathValue("itemId"))
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

	cartID, err := strconv.Atoi(r.PathValue("id"))
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

// GET http://localhost:3000/carts/1/price
func (h *Handler) getCartPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	cartID, err := strconv.Atoi(r.PathValue("id"))
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
