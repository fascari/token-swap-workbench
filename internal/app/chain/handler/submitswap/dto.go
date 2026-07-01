package submitswap

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
