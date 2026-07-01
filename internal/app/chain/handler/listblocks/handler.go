package listblocks

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
	"github.com/fascari/token-swap-workbench/pkg/httpparam"
)

const (
	path       = "/blocks"
	queryCount = "n"
)

type (
	Handler struct {
		useCase listblocks.UseCase
	}
)

func New(useCase listblocks.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterRoutes(r chi.Router, h Handler) {
	r.Get(path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	count, err := httpparam.PositiveInt(r.URL.Query().Get(queryCount), queryCount)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	output, err := h.useCase.Execute(r.Context(), listblocks.Input{Count: count})
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, domain.ErrUpstreamRejected) {
			status = http.StatusBadRequest
		}
		httpjson.WriteError(w, status, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, toResponse(output.Blocks))
}
