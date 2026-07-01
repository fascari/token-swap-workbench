package validator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/pkg/validator"
	validatortestdata "github.com/fascari/token-swap-workbench/pkg/validator/testdata"
)

func TestValidate_ShouldReturnNilWhenStructIsValid(t *testing.T) {
	err := validator.Validate(validatortestdata.ValidPayload())

	require.NoError(t, err)
}

func TestValidate_ShouldReturnErrorWhenStructIsInvalid(t *testing.T) {
	tests := []struct {
		name          string
		input         validatortestdata.Payload
		expectedError string
	}{
		{
			name:          "should return error when account_id is zero",
			input:         validatortestdata.PayloadWithZeroAccountID(),
			expectedError: "account_id must be greater than 0",
		},
		{
			name:          "should return error when in_token is empty",
			input:         validatortestdata.PayloadWithEmptyInToken(),
			expectedError: "in_token is a required field",
		},
		{
			name:          "should return error when all fields are invalid",
			input:         validatortestdata.InvalidPayload(),
			expectedError: "account_id must be greater than 0; in_token is a required field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.input)

			require.EqualError(t, err, tt.expectedError)
		})
	}
}
