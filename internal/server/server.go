package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/logger"
	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/server/config"
)

type Server struct {
	httpServer *http.Server
	router     http.Handler
	logger     logger.Logger
}

func New(logger logger.Logger, handler http.Handler, config *config.Config) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         config.Address(),
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context) error {

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	s.logger.Info("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(shutdownCtx)
}
