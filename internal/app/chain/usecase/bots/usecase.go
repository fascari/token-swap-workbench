package bots

import (
	"context"
	"fmt"
)

const (
	maxBotsPerRequest = 100
	statusRunning     = "running"
	statusStopped     = "stopped"
)

func New(client Client) UseCase {
	return UseCase{
		client:  client,
		manager: newManager(),
	}
}

func (uc UseCase) Execute(ctx context.Context, input Input) (Output, error) {
	if err := input.validate(); err != nil {
		return Output{}, err
	}

	switch input.Action {
	case ActionCreate:
		created := uc.manager.create(ctx, uc.client, input.Amount)
		return uc.manager.output(input, created, 0), nil
	case ActionStop:
		if input.All {
			stopped := uc.manager.stopAll()
			input.Amount = stopped
			return uc.manager.output(input, 0, stopped), nil
		}
		stopped := uc.manager.stop(input.Amount)
		return uc.manager.output(input, 0, stopped), nil
	default:
		return Output{}, fmt.Errorf("%w: %s", ErrInvalidAction, input.Action)
	}
}
