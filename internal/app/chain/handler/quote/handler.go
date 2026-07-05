package quote

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
	"github.com/fascari/token-swap-workbench/pkg/httpjson"
	"github.com/fascari/token-swap-workbench/pkg/httpparam"
)

const (
	path          = "/quote"
	queryInToken  = "in"
	queryOutToken = "out"
	queryAmount   = "amount"
)

type (
	Handler struct {
		useCase quote.UseCase
	}
)

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

	amount, err := httpparam.PositiveFloat(r.URL.Query().Get(queryAmount), queryAmount)
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
		httpjson.WriteError(w, errorCode(err), err)
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, responseDTO{AmountOut: output.AmountOut})
}

func token(value, field string) (domain.Token, error) {
	raw, err := httpparam.Required(value, field)
	if err != nil {
		return "", err
	}

	return domain.Token(raw), nil
}
