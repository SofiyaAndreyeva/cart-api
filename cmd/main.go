package main

import (
	"log/slog"

	"github.com/SofiyaAndreyeva/cart-api/internal/config"
	"github.com/SofiyaAndreyeva/cart-api/internal/handler"
	"github.com/SofiyaAndreyeva/cart-api/internal/repository"
	"github.com/SofiyaAndreyeva/cart-api/internal/server"
	"github.com/SofiyaAndreyeva/cart-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("error initializing configs", "error", err)
		return
	}

	repo := repository.NewRepository()
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	serv := server.NewServer()
	if err := serv.Run(cfg.HTTPPort, handlers.InitialRoutes()); err != nil {
		slog.Error("application stopped", "error", err)
	}
}
