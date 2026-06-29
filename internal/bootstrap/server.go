package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fascari/token-swap-workbench/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// App holds the application state and dependencies.
type App struct {
	config *config.Config

	router chi.Router
	server *http.Server
}

// New creates and initialises a new application instance.
func New(ctx context.Context, cfg *config.Config) (*App, func(), error) {
	app := &App{config: cfg}

	cleanupFuncs := []func(){}

	app.router = NewRouter()

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: app.router,
	}

	cleanup := func() {
		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}
	}

	return app, cleanup, nil
}

// Run starts the HTTP server and blocks until the context is cancelled.
func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		log.Info().Int("port", a.config.HTTP.Port).Msg("starting HTTP server")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("http server error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("shutting down gracefully")
	case err := <-errChan:
		log.Error().Err(err).Msg("server error")
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("http server shutdown error")
		return err
	}

	log.Info().Msg("shutdown complete")
	return nil
}
