package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fascari/token-swap-workbench/cmd/api/modules"
	"github.com/fascari/token-swap-workbench/internal/chainclient"
	"github.com/fascari/token-swap-workbench/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type App struct {
	router chi.Router
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	chainClient, err := chainclient.New(cfg.Chain)
	if err != nil {
		return nil, fmt.Errorf("creating chain client: %w", err)
	}

	app := &App{}
	app.router = NewRouter(modules.NewChainModule(chainClient))

	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: app.router,
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		log.Info().Str("addr", a.server.Addr).Msg("starting HTTP server")
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
