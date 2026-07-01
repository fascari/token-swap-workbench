package quote_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote/mocks"
)

func TestUseCase_Execute_ShouldReturnQuoteWhenClientQuotes(t *testing.T) {
	request := chaintestdata.QuoteRequest()
	expectedQuote := chaintestdata.Quote()
	client := mocks.NewClient(t)
	client.EXPECT().Quote(t.Context(), request).Return(expectedQuote, nil)

	uc := quote.NewUseCase(client)

	output, err := uc.Execute(t.Context(), quote.Input{
		InToken:  request.InToken,
		OutToken: request.OutToken,
		Amount:   request.Amount,
	})

	require.NoError(t, err)
	require.Equal(t, quote.Output{AmountOut: expectedQuote.AmountOut}, output)
}

func TestUseCase_Execute_ShouldReturnErrorWhenClientQuoteFails(t *testing.T) {
	request := chaintestdata.QuoteRequest()
	expectedErr := errors.New("quote unavailable")
	client := mocks.NewClient(t)
	client.EXPECT().Quote(t.Context(), request).Return(chaintestdata.Quote(), expectedErr)

	uc := quote.NewUseCase(client)

	output, err := uc.Execute(t.Context(), quote.Input{
		InToken:  request.InToken,
		OutToken: request.OutToken,
		Amount:   request.Amount,
	})

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, quote.Output{}, output)
}
