package bots

import (
	"context"
	"errors"
	"fmt"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

const (
	ActionCreate Action = "create"
	ActionStop   Action = "stop"
)

var (
	ErrInvalidAction = errors.New("invalid bot action")
	ErrInvalidAmount = errors.New("invalid bot amount")
)

type (
	Action string

	Client interface {
		SubmitTransaction(ctx context.Context, transaction domain.TransactionSubmission) error
	}

	UseCase struct {
		client  Client
		manager *manager
	}

	Input struct {
		Action Action
		Amount int
		All    bool
	}

	Output struct {
		Status              string
		Action              Action
		RequestedAmount     int
		All                 bool
		ActiveBots          int
		CreatedBots         int
		StoppedBots         int
		AttemptedOperations int
		AcceptedOperations  int
		FailedOperations    int
		SendOperations      int
		SwapOperations      int
	}
)

func (input Input) validate() error {
	switch input.Action {
	case ActionCreate, ActionStop:
	default:
		return fmt.Errorf("%w: %s", ErrInvalidAction, input.Action)
	}

	if input.All && input.Action != ActionStop {
		return fmt.Errorf("%w: all can only be used with %s", ErrInvalidAction, ActionStop)
	}

	if input.All {
		return nil
	}

	if input.Amount <= 0 || input.Amount > maxBotsPerRequest {
		return fmt.Errorf("%w: %d", ErrInvalidAmount, input.Amount)
	}

	return nil
}
