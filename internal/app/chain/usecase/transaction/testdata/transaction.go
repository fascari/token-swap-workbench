package testdata

import (
	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction"
)

func SendSubmission() domain.TransactionSubmission {
	return domain.TransactionSubmission{
		Kind: domain.TransactionKindSend,
		Send: chaintestdata.Send(),
	}
}

func SwapSubmission() domain.TransactionSubmission {
	return domain.TransactionSubmission{
		Kind: domain.TransactionKindSwap,
		Swap: chaintestdata.Swap(),
	}
}

func Output() transaction.Output {
	return transaction.Output{Status: domain.StatusSubmitted}
}
