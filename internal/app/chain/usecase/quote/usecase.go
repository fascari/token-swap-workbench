package quote

import (
	"context"
	"fmt"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

type (
	Client interface {
		Quote(ctx context.Context, req domain.QuoteRequest) (domain.Quote, error)
	}

	UseCase struct {
		client Client
	}

	Input struct {
		InToken  domain.Token
		OutToken domain.Token
		Amount   float64
	}

	Output struct {
		AmountOut float64
	}
)

func New(client Client) UseCase {
	return UseCase{client: client}
}

func (uc UseCase) Execute(ctx context.Context, input Input) (Output, error) {
	quote, err := uc.client.Quote(ctx, domain.QuoteRequest{
		InToken:  input.InToken,
		OutToken: input.OutToken,
		Amount:   input.Amount,
	})
	if err != nil {
		return Output{}, fmt.Errorf("quoting swap: %w", err)
	}

	return Output{AmountOut: quote.AmountOut}, nil
}
