package testdata

import (
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

func QuoteRequest() domain.QuoteRequest {
	return chaintestdata.QuoteRequest()
}

func Quote() domain.Quote {
	return chaintestdata.Quote()
}
