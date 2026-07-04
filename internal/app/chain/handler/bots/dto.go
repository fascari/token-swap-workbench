package bots

import botsuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots"

type (
	requestDTO struct {
		Action string `json:"action" validate:"required"`
		Amount int    `json:"amount"`
		All    bool   `json:"all"`
	}

	responseDTO struct {
		Status              string `json:"status"`
		Action              string `json:"action"`
		RequestedAmount     int    `json:"requested_amount"`
		All                 bool   `json:"all"`
		ActiveBots          int    `json:"active_bots"`
		CreatedBots         int    `json:"created_bots"`
		StoppedBots         int    `json:"stopped_bots"`
		AttemptedOperations int    `json:"attempted_operations"`
		AcceptedOperations  int    `json:"accepted_operations"`
		FailedOperations    int    `json:"failed_operations"`
		SendOperations      int    `json:"send_operations"`
		SwapOperations      int    `json:"swap_operations"`
	}
)

func (dto requestDTO) toDomain() botsuc.Input {
	return botsuc.Input{
		Action: botsuc.Action(dto.Action),
		Amount: dto.Amount,
		All:    dto.All,
	}
}

func toResponse(output botsuc.Output) responseDTO {
	return responseDTO{
		Status:              output.Status,
		Action:              string(output.Action),
		RequestedAmount:     output.RequestedAmount,
		All:                 output.All,
		ActiveBots:          output.ActiveBots,
		CreatedBots:         output.CreatedBots,
		StoppedBots:         output.StoppedBots,
		AttemptedOperations: output.AttemptedOperations,
		AcceptedOperations:  output.AcceptedOperations,
		FailedOperations:    output.FailedOperations,
		SendOperations:      output.SendOperations,
		SwapOperations:      output.SwapOperations,
	}
}
