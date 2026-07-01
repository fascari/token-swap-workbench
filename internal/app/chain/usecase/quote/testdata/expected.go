package testdata

import (
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
)

func Output() quote.Output {
	return quote.Output{AmountOut: chaintestdata.QuoteAmountOut}
}
