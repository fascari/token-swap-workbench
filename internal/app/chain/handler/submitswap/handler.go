package submitswap

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
	"github.com/fascari/token-swap-workbench/pkg/validator"
)

const (
	path = "/swaps"
)

type (
	Handler struct {
		useCase submitswap.UseCase
	}
)

func New(useCase submitswap.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterRoutes(r chi.Router, h Handler) {
	r.Post(path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload requestDTO
	if err := httpjson.ReadJSON(r, &payload); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := validator.Validate(payload); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	output, err := h.useCase.Execute(r.Context(), submitswap.Input{
		AccountID: payload.AccountID,
		InToken:   domain.Token(payload.InToken),
		OutToken:  domain.Token(payload.OutToken),
		AmountIn:  payload.AmountIn,
	})
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, domain.ErrUpstreamRejected) {
			status = http.StatusBadRequest
		}
		httpjson.WriteError(w, status, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusAccepted, responseDTO{Status: output.Status})
}
