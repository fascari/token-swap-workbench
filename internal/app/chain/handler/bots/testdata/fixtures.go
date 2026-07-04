package testdata

import _ "embed"

var (
	//go:embed json/create_request.json
	CreateRequest string

	//go:embed json/create_response.json
	CreateResponse string

	//go:embed json/stop_all_request.json
	StopAllRequest string

	//go:embed json/stop_all_response.json
	StopAllResponse string

	//go:embed json/malformed_request.json
	MalformedRequest string

	//go:embed json/missing_action_request.json
	MissingActionRequest string

	//go:embed json/invalid_amount_request.json
	InvalidAmountRequest string
)
