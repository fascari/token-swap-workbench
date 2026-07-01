package testdata

import (
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
)

func Output() listblocks.Output {
	return listblocks.Output{Blocks: chaintestdata.Blocks()}
}
