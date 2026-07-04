package testdata

import "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots"

const (
	ValidAmount    = 2
	TooLargeAmount = 101
)

func CreateInput() bots.Input {
	return bots.Input{
		Action: bots.ActionCreate,
		Amount: ValidAmount,
	}
}

func StopInput() bots.Input {
	return bots.Input{
		Action: bots.ActionStop,
		Amount: ValidAmount,
	}
}

func StopAllInput() bots.Input {
	return bots.Input{
		Action: bots.ActionStop,
		All:    true,
	}
}

func InvalidActionInput() bots.Input {
	return bots.Input{Action: "invalid"}
}

func AllWithCreateInput() bots.Input {
	return bots.Input{
		Action: bots.ActionCreate,
		All:    true,
	}
}

func ZeroAmountInput() bots.Input {
	return bots.Input{
		Action: bots.ActionCreate,
		Amount: 0,
	}
}

func TooLargeAmountInput() bots.Input {
	return bots.Input{
		Action: bots.ActionCreate,
		Amount: TooLargeAmount,
	}
}

func CreatedOutput() bots.Output {
	return bots.Output{
		Status:          "running",
		Action:          bots.ActionCreate,
		RequestedAmount: ValidAmount,
		ActiveBots:      ValidAmount,
		CreatedBots:     ValidAmount,
	}
}

func StoppedOutput() bots.Output {
	return bots.Output{
		Status:          "stopped",
		Action:          bots.ActionStop,
		RequestedAmount: ValidAmount,
		StoppedBots:     ValidAmount,
	}
}

func StoppedAllOutput() bots.Output {
	return bots.Output{
		Status:          "stopped",
		Action:          bots.ActionStop,
		RequestedAmount: ValidAmount,
		All:             true,
		StoppedBots:     ValidAmount,
	}
}
