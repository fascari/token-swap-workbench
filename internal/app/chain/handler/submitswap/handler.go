package submitswap

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
)

const path = "/swaps"

type Handler struct {
	useCase submitswap.UseCase
}

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

	if err := validate(payload); err != nil {
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
		httpjson.WriteError(w, http.StatusBadGateway, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusAccepted, responseDTO{Status: output.Status})
}

func validate(payload requestDTO) error {
	if payload.AccountID == 0 {
		return errors.New("account_id must be greater than 0")
	}

	if payload.InToken == "" {
		return errors.New("in_token is required")
	}

	if payload.OutToken == "" {
		return errors.New("out_token is required")
	}

	if payload.AmountIn <= 0 {
		return fmt.Errorf("%s must be greater than 0", "amount_in")
	}

	return nil
}
