package testdata

import (
	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status"
)

func Output() status.Output {
	return status.Output{Status: domain.StatusOK}
}
