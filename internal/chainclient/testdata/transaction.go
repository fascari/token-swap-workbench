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

func Swap() domain.Swap {
	return domain.Swap{
		AccountID: AccountID,
		InToken:   TokenUSDC,
		OutToken:  TokenETH,
		AmountIn:  AmountIn,
	}
}

func Send() domain.Send {
	return domain.Send{
		From:   AccountID,
		To:     RecipientID,
		Amount: SendAmount,
		Token:  TokenUSDC,
	}
}

func SwapSubmission() domain.TransactionSubmission {
	return domain.TransactionSubmission{
		Kind: domain.TransactionKindSwap,
		Swap: Swap(),
	}
}

func SendSubmission() domain.TransactionSubmission {
	return domain.TransactionSubmission{
		Kind: domain.TransactionKindSend,
		Send: Send(),
	}
}
