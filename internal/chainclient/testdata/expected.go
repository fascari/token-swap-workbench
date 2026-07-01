package testdata

import "github.com/fascari/token-swap-workbench/internal/app/chain/domain"

func Quote() domain.Quote {
	return domain.Quote{AmountOut: AmountOut}
}

func Blocks() []domain.Block {
	return []domain.Block{
		{
			ID:        BlockID,
			Timestamp: BlockTime,
			Transactions: []domain.Transaction{
				{
					Swap: new(domain.SwapTransaction{
						AccountID: AccountID,
						InToken:   TokenUSDC,
						OutToken:  TokenETH,
						AmountIn:  AmountIn,
					}),
				},
				{
					Send: new(domain.SendTransaction{
						From:   AccountID,
						To:     RecipientID,
						Amount: SendAmount,
						Token:  TokenUSDC,
					}),
				},
			},
		},
	}
}
