package testdata

import _ "embed"

var (
	//go:embed json/quote_response.json
	QuoteResponse string

	//go:embed json/blocks_response.json
	BlocksResponse string

	//go:embed json/send_envelope.json
	SendEnvelope string

	//go:embed json/swap_envelope.json
	SwapEnvelope string
)
