package testdata

import (
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

func Blocks() []domain.Block {
	return chaintestdata.Blocks()
}
