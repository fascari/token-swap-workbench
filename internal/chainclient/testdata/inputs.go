package testdata

import "github.com/fascari/token-swap-workbench/internal/app/chain/domain"

const (
	AccountID   uint32       = 42
	RecipientID uint32       = 84
	AmountIn    float64      = 12.5
	SendAmount  float64      = 3.75
	AmountOut   float64      = 6.25
	BlockID     uint64       = 99
	BlockTime   uint64       = 1717000000
	TokenUSDC   domain.Token = "USDC"
	TokenETH    domain.Token = "ETH"
)

func QuoteRequest() domain.QuoteRequest {
	return domain.QuoteRequest{
		InToken:  TokenUSDC,
		OutToken: TokenETH,
		Amount:   AmountIn,
	}
}

func Swap() domain.Swap {
	return domain.Swap{
		AccountID: AccountID,
		InToken:   TokenUSDC,
		OutToken:  TokenETH,
		AmountIn:  AmountIn,
	}
}
