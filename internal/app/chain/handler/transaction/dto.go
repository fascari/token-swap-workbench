package transaction

import "github.com/fascari/token-swap-workbench/internal/app/chain/domain"

type (
	requestDTO struct {
		AccountID uint32  `json:"account_id" validate:"gt=0"`
		InToken   string  `json:"in_token" validate:"required"`
		OutToken  string  `json:"out_token" validate:"required"`
		AmountIn  float64 `json:"amount_in" validate:"gt=0"`
	}

	responseDTO struct {
		Status string `json:"status"`
	}
)

func (dto requestDTO) toDomain() domain.TransactionSubmission {
	return domain.TransactionSubmission{
		Kind: domain.TransactionKindSwap,
		Swap: domain.Swap{
			AccountID: dto.AccountID,
			InToken:   domain.Token(dto.InToken),
			OutToken:  domain.Token(dto.OutToken),
			AmountIn:  dto.AmountIn,
		},
	}
}
