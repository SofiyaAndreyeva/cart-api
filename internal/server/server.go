package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SofiyaAndreyeva/cart-api/internal/domain"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{
		httpServer: &http.Server{
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer.Addr = ":" + port
	s.httpServer.Handler = handler
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrServerRun, err)
	}

	return nil
}
