package handler

import (
	"errors"
	"net/http"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
)

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCartNotFound), errors.Is(err, domain.ErrCartItemNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrLimitCart):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, domain.ErrEmptyProduct), errors.Is(err, domain.ErrInvalidPrice):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
