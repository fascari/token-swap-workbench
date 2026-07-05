package testdata

import (
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

var (
	AmountIn  = chaintestdata.SwapAmountIn
	AmountOut = chaintestdata.QuoteAmountOut
)

func SwapSubmission() domain.TransactionSubmission {
	return domain.TransactionSubmission{
		Kind: domain.TransactionKindSwap,
		Swap: chaintestdata.Swap(),
	}
}

func SendSubmission() domain.TransactionSubmission {
	return domain.TransactionSubmission{
		Kind: domain.TransactionKindSend,
		Send: chaintestdata.Send(),
	}
}
