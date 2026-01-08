package main

import (
	"fmt"
	"log/slog"

	"github.com/SofiyaAndreyeva/cart-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("error initializing configs", "error", err)
		return
	}
	fmt.Println(cfg)
}
