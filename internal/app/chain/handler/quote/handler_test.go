package quote_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/quote"
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	quoteuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote/mocks"
)

func TestHandler_Handle_ShouldReturnQuoteWhenRequestIsValid(t *testing.T) {
	quoteRequest := chaintestdata.QuoteRequest()
	expectedQuote := chaintestdata.Quote()
	client := mocks.NewClient(t)
	client.EXPECT().Quote(t.Context(), quoteRequest).Return(expectedQuote, nil)
	handler := quote.New(quoteuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		fmt.Sprintf(
			"/chain/quote?in=%s&out=%s&amount=%.1f",
			quoteRequest.InToken,
			quoteRequest.OutToken,
			quoteRequest.Amount,
		),
		nil,
	)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, `{"amount_out":6.25}`, recorder.Body.String())
}

func TestHandler_Handle_ShouldReturnBadRequestWhenAmountIsMissing(t *testing.T) {
	client := mocks.NewClient(t)
	handler := quote.New(quoteuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/chain/quote?in=NEX&out=ETH", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Handle_ShouldReturnBadGatewayWhenQuoteFails(t *testing.T) {
	quoteRequest := chaintestdata.QuoteRequest()
	client := mocks.NewClient(t)
	client.EXPECT().Quote(t.Context(), quoteRequest).Return(chaintestdata.Quote(), errors.New("quote unavailable"))
	handler := quote.New(quoteuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		fmt.Sprintf(
			"/chain/quote?in=%s&out=%s&amount=%.1f",
			quoteRequest.InToken,
			quoteRequest.OutToken,
			quoteRequest.Amount,
		),
		nil,
	)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadGateway, recorder.Code)
}
