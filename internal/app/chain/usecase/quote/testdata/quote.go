package testdata

import (
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
)

func Input() quote.Input {
	req := chaintestdata.QuoteRequest()
	return quote.Input{
		InToken:  req.InToken,
		OutToken: req.OutToken,
		Amount:   req.Amount,
	}
}

func Output() quote.Output {
	return quote.Output{AmountOut: chaintestdata.QuoteAmountOut}
}
