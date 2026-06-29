package bootstrap

import (
	"net/http"

	"github.com/fascari/token-swap-workbench/cmd/api/modules"
	"github.com/fascari/token-swap-workbench/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates and configures the HTTP router with all middleware and domain modules.

func NewRouter() chi.Router {
	exampleModule := modules.NewExampleModule()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger())
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index.html")
	})
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.Route("/v1", func(r chi.Router) {
		exampleModule.Register(r)
	})

	return r
}
