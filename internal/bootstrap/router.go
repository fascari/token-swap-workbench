package bootstrap

import (
	"encoding/json"
	"net/http"

	"github.com/fascari/token-swap-workbench/cmd/api/modules"
	"github.com/fascari/token-swap-workbench/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(modulesList ...modules.Module) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger())
	r.Use(chimiddleware.Recoverer)

	registerHealthRoute(r)
	registerFrontendRoutes(r)

	r.Route("/v1", func(r chi.Router) {
		for _, module := range modulesList {
			module.Register(r)
		}
	})

	return r
}

func registerHealthRoute(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})
}

func registerFrontendRoutes(r chi.Router) {
	r.Get("/", serveIndex)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/index.html")
}
