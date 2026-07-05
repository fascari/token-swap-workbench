package transaction

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	transactionuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
	"github.com/fascari/token-swap-workbench/pkg/validator"
)

const (
	path = "/transactions"
)

type (
	Handler struct {
		useCase transactionuc.UseCase
	}
)

func New(useCase transactionuc.UseCase) Handler {
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

	output, err := h.useCase.Execute(r.Context(), payload.toDomain())
	if err != nil {
		httpjson.WriteError(w, errorCode(err), err)
		return
	}

	httpjson.WriteJSON(w, http.StatusAccepted, responseDTO{Status: output.Status})
}
