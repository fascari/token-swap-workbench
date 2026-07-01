package submitswap

type (
	requestDTO struct {
		AccountID uint32  `json:"account_id"`
		InToken   string  `json:"in_token"`
		OutToken  string  `json:"out_token"`
		AmountIn  float64 `json:"amount_in"`
	}

	responseDTO struct {
		Status string `json:"status"`
	}
)
