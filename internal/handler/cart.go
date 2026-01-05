package handler

import (
	"encoding/json"
	"net/http"
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
