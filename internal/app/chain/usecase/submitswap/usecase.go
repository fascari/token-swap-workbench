package submitswap

import (
	"context"
	"fmt"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

type (
	Client interface {
		SubmitSwap(ctx context.Context, swap domain.Swap) error
	}

	UseCase struct {
		client Client
	}

	Input struct {
		AccountID uint32
		InToken   domain.Token
		OutToken  domain.Token
		AmountIn  float64
	}

	Output struct {
		Status string
	}
)

func NewUseCase(client Client) UseCase {
	return UseCase{client: client}
}

func (uc UseCase) Execute(ctx context.Context, input Input) (Output, error) {
	err := uc.client.SubmitSwap(ctx, domain.Swap{
		AccountID: input.AccountID,
		InToken:   input.InToken,
		OutToken:  input.OutToken,
		AmountIn:  input.AmountIn,
	})
	if err != nil {
		return Output{}, fmt.Errorf("submitting swap: %w", err)
	}

	return Output{Status: "submitted"}, nil
}
