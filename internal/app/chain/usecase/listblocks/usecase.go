package listblocks

import (
	"context"
	"fmt"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

type (
	Client interface {
		Blocks(ctx context.Context, n int) ([]domain.Block, error)
	}

	UseCase struct {
		client Client
	}

	Input struct {
		Count int
	}

	Output struct {
		Blocks []domain.Block
	}
)

func New(client Client) UseCase {
	return UseCase{client: client}
}

func (uc UseCase) Execute(ctx context.Context, input Input) (Output, error) {
	blocks, err := uc.client.Blocks(ctx, input.Count)
	if err != nil {
		return Output{}, fmt.Errorf("listing blocks: %w", err)
	}

	return Output{Blocks: blocks}, nil
}
