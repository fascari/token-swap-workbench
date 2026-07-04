package testdata

import _ "embed"

var (
	//go:embed json/transaction_request.json
	Request string

	//go:embed json/invalid_transaction_request.json
	InvalidRequest string

	//go:embed json/response.json
	Response string
)
