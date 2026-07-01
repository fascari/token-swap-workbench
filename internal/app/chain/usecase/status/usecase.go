package status

import (
	"context"
	"fmt"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

type (
	Client interface {
		Status(ctx context.Context) error
	}

	UseCase struct {
		client Client
	}

	Input struct{}

	Output struct {
		Status string
	}
)

func NewUseCase(client Client) UseCase {
	return UseCase{client: client}
}

func (uc UseCase) Execute(ctx context.Context, _ Input) (Output, error) {
	if err := uc.client.Status(ctx); err != nil {
		return Output{}, fmt.Errorf("checking chain status: %w", err)
	}

	return Output{Status: domain.StatusOK}, nil
}
