package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/RAdevelop/ya_practicum-go-advanced-ext-diplom/internal/app_context"
)

type Server struct {
	httpServer *http.Server
	appContext *app_context.AppContext
}

func New(appContext *app_context.AppContext, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         appContext.ServerConfig.Address(),
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		appContext: appContext,
	}
}

func (s *Server) Run(ctx context.Context) error {

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.appContext.Logger.Error("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	s.appContext.Logger.Info("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(shutdownCtx)
}
