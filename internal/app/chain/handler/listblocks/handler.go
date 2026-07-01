package listblocks

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
)

const (
	path       = "/blocks"
	queryCount = "n"
)

type Handler struct {
	useCase listblocks.UseCase
}

func New(useCase listblocks.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterRoutes(r chi.Router, h Handler) {
	r.Get(path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	count, err := positiveInt(r.URL.Query().Get(queryCount), queryCount)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	output, err := h.useCase.Execute(r.Context(), listblocks.Input{Count: count})
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, toResponse(output.Blocks))
}

func positiveInt(value string, field string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("%s is required", field)
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", field, err)
	}

	if number <= 0 {
		return 0, errors.New(field + " must be greater than 0")
	}

	return number, nil
}
