package transaction

import (
	"context"
	"fmt"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

type (
	Client interface {
		SubmitTransaction(ctx context.Context, transaction domain.TransactionSubmission) error
	}

	UseCase struct {
		client Client
	}

	Output struct {
		Status string
	}
)

func New(client Client) UseCase {
	return UseCase{client: client}
}

func (uc UseCase) Execute(ctx context.Context, submission domain.TransactionSubmission) (Output, error) {
	err := uc.client.SubmitTransaction(ctx, submission)
	if err != nil {
		return Output{}, fmt.Errorf("submitting transaction: %w", err)
	}

	return Output{Status: domain.StatusSubmitted}, nil
}
