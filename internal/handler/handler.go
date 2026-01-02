package handler

import (
	"net/http"

	"github.com/SofiyaAndreyeva/cart-api/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitialRoutes() http.Handler {
	mux := http.NewServeMux()
	return mux
}
