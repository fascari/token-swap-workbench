package domain

const (
	StatusOK        = "ok"
	StatusSubmitted = "submitted"
)

type (
	Token string

	QuoteRequest struct {
		InToken  Token
		OutToken Token
		Amount   float64
	}

	Quote struct {
		AmountOut float64
	}

	Swap struct {
		AccountID uint32
		InToken   Token
		OutToken  Token
		AmountIn  float64
	}

	Block struct {
		ID           uint64
		Timestamp    uint64
		Transactions []Transaction
	}

	Transaction struct {
		Swap *SwapTransaction
		Send *SendTransaction
	}

	SwapTransaction struct {
		AccountID uint32
		InToken   Token
		OutToken  Token
		AmountIn  float64
	}

	SendTransaction struct {
		From   uint32
		To     uint32
		Amount float64
		Token  Token
	}
)
