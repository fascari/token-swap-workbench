package quote

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
)

const (
	path          = "/quote"
	queryInToken  = "in"
	queryOutToken = "out"
	queryAmount   = "amount"
)

type Handler struct {
	useCase quote.UseCase
}

func New(useCase quote.UseCase) Handler {
	return Handler{useCase: useCase}
}

func RegisterRoutes(r chi.Router, h Handler) {
	r.Get(path, h.Handle)
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	inToken, err := token(r.URL.Query().Get(queryInToken), queryInToken)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	outToken, err := token(r.URL.Query().Get(queryOutToken), queryOutToken)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	amount, err := positiveFloat(r.URL.Query().Get(queryAmount), queryAmount)
	if err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	output, err := h.useCase.Execute(r.Context(), quote.Input{
		InToken:  inToken,
		OutToken: outToken,
		Amount:   amount,
	})
	if err != nil {
		httpjson.WriteError(w, http.StatusBadGateway, err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, responseDTO{AmountOut: output.AmountOut})
}

func token(value string, field string) (domain.Token, error) {
	if value == "" {
		return "", fmt.Errorf("%s is required", field)
	}

	return domain.Token(value), nil
}

func positiveFloat(value string, field string) (float64, error) {
	if value == "" {
		return 0, fmt.Errorf("%s is required", field)
	}

	number, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number: %w", field, err)
	}

	if number <= 0 {
		return 0, errors.New(field + " must be greater than 0")
	}

	return number, nil
}
