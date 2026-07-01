package testdata

import "github.com/fascari/token-swap-workbench/internal/app/chain/domain"

const (
	AccountID           uint32 = 42
	RecipientID         uint32 = 84
	BlockID             uint64 = 99
	BlockTime           uint64 = 1717000000
	TokenUSDC                  = domain.Token("USDC")
	TokenETH                   = domain.Token("ETH")
	SwapStatusSubmitted        = "submitted"
	SwapAmountIn               = 12.5
	QuoteAmountOut             = 6.25
	SendAmount                 = 3.75
)

func QuoteRequest() domain.QuoteRequest {
	return domain.QuoteRequest{
		InToken:  TokenUSDC,
		OutToken: TokenETH,
		Amount:   SwapAmountIn,
	}
}

func Quote() domain.Quote {
	return domain.Quote{
		AmountOut: QuoteAmountOut,
	}
}

func Swap() domain.Swap {
	return domain.Swap{
		AccountID: AccountID,
		InToken:   TokenUSDC,
		OutToken:  TokenETH,
		AmountIn:  SwapAmountIn,
	}
}

func Blocks() []domain.Block {
	return []domain.Block{
		{
			ID:        BlockID,
			Timestamp: BlockTime,
			Transactions: []domain.Transaction{
				{
					Swap: &domain.SwapTransaction{
						AccountID: AccountID,
						InToken:   TokenUSDC,
						OutToken:  TokenETH,
						AmountIn:  SwapAmountIn,
					},
				},
				{
					Send: &domain.SendTransaction{
						From:   AccountID,
						To:     RecipientID,
						Amount: SendAmount,
						Token:  TokenUSDC,
					},
				},
			},
		},
	}
}
