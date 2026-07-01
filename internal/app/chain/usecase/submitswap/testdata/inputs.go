package testdata

import (
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
)

func Input() submitswap.Input {
	swap := chaintestdata.Swap()
	return submitswap.Input{
		AccountID: swap.AccountID,
		InToken:   swap.InToken,
		OutToken:  swap.OutToken,
		AmountIn:  swap.AmountIn,
	}
}
