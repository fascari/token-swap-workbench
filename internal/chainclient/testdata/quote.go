package testdata

import "github.com/fascari/token-swap-workbench/internal/app/chain/domain"

func QuoteRequest() domain.QuoteRequest {
	return domain.QuoteRequest{
		InToken:  TokenUSDC,
		OutToken: TokenETH,
		Amount:   AmountIn,
	}
}

func Quote() domain.Quote {
	return domain.Quote{AmountOut: AmountOut}
}
