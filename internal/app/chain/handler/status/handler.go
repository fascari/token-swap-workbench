package status

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
)

const (
	path = "/chain/status"
)

type (
	Handler struct {
		useCase status.UseCase
	}
)

func New(useCase status.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterRoutes(r chi.Router, h Handler) {
	r.Get(path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	output, err := h.useCase.Execute(r.Context(), status.Input{})
	if err != nil {
		httpjson.WriteError(w, errorCode(err), err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, responseDTO{Status: output.Status})
}
