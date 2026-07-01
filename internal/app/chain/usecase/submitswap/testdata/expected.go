package testdata

import (
	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
)

func Output() submitswap.Output {
	return submitswap.Output{Status: domain.StatusSubmitted}
}
