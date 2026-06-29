package modules

import (
	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/example/handler/createexample"
	examplerepo "github.com/fascari/token-swap-workbench/internal/app/example/repository"
	createexampleuc "github.com/fascari/token-swap-workbench/internal/app/example/usecase/createexample"
)

// ExampleModule wires the example domain: repository → usecase → handler.
type ExampleModule struct {
	handler createexample.Handler
}

func NewExampleModule() *ExampleModule {
	repo := examplerepo.New()

	uc := createexampleuc.New(repo)
	return &ExampleModule{handler: createexample.NewHandler(uc)}
}

func (m *ExampleModule) Register(r chi.Router) {
	createexample.RegisterRoutes(r, m.handler)
}
