package chainclient

const (
	TokenUSDC Token = "USDC"
	TokenUSDT Token = "USDT"
	TokenETH  Token = "ETH"
	TokenBTC  Token = "BTC"
	TokenNEX  Token = "NEX"
	TokenDOGE Token = "DOGE"
	TokenHYPE Token = "HYPE"
)

type (
	Token string

	QuoteRequest struct {
		InToken  Token
		OutToken Token
		Amount   float64
	}

	QuoteResponse struct {
		AmountOut float64 `json:"amount_out"`
	}

	SwapRequest struct {
		AccountID uint32  `json:"account"`
		InToken   Token   `json:"in_token"`
		OutToken  Token   `json:"out_token"`
		AmountIn  float64 `json:"amount_in"`
	}

	Block struct {
		ID           uint64        `json:"id"`
		Timestamp    uint64        `json:"timestamp"`
		Transactions []Transaction `json:"transactions"`
	}

	Transaction struct {
		Swap *SwapTransaction `json:"Swap,omitzero"`
		Send *SendTransaction `json:"Send,omitzero"`
	}

	SwapTransaction struct {
		AccountID uint32  `json:"account_id"`
		InToken   Token   `json:"in_token"`
		OutToken  Token   `json:"out_token"`
		AmountIn  float64 `json:"amount_in"`
	}

	SendTransaction struct {
		From   uint32  `json:"from"`
		To     uint32  `json:"to"`
		Amount float64 `json:"amount"`
		Token  Token   `json:"token"`
	}
)
