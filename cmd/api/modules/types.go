package modules

import "github.com/go-chi/chi/v5"

type Module interface {
	Register(r chi.Router)
}
