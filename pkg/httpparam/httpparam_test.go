package httpparam_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/pkg/httpparam"
)

func TestRequired_ShouldReturnValueWhenPresent(t *testing.T) {
	value, err := httpparam.Required("NEX", "in")

	require.NoError(t, err)
	require.Equal(t, "NEX", value)
}

func TestRequired_ShouldReturnErrorWhenEmpty(t *testing.T) {
	_, err := httpparam.Required("", "in")

	require.EqualError(t, err, "in is required")
}

func TestPositiveInt_ShouldReturnValueWhenPositive(t *testing.T) {
	value, err := httpparam.PositiveInt("10", "n")

	require.NoError(t, err)
	require.Equal(t, 10, value)
}

func TestPositiveInt_ShouldReturnErrorWhenInvalid(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		expectedError string
	}{
		{
			name:          "should return error when empty",
			value:         "",
			expectedError: "n is required",
		},
		{
			name:          "should return error when value is not an integer",
			value:         "abc",
			expectedError: `n must be an integer: strconv.Atoi: parsing "abc": invalid syntax`,
		},
		{
			name:          "should return error when zero",
			value:         "0",
			expectedError: "n must be greater than 0",
		},
		{
			name:          "should return error when negative",
			value:         "-3",
			expectedError: "n must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := httpparam.PositiveInt(tt.value, "n")

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestPositiveFloat_ShouldReturnValueWhenPositive(t *testing.T) {
	value, err := httpparam.PositiveFloat("12.5", "amount")

	require.NoError(t, err)
	require.InEpsilon(t, 12.5, value, 1e-9)
}

func TestPositiveFloat_ShouldReturnErrorWhenInvalid(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		expectedError string
	}{
		{
			name:          "should return error when empty",
			value:         "",
			expectedError: "amount is required",
		},
		{
			name:          "should return error when value is not a number",
			value:         "abc",
			expectedError: `amount must be a number: strconv.ParseFloat: parsing "abc": invalid syntax`,
		},
		{
			name:          "should return error when zero",
			value:         "0",
			expectedError: "amount must be greater than 0",
		},
		{
			name:          "should return error when negative",
			value:         "-1.5",
			expectedError: "amount must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := httpparam.PositiveFloat(tt.value, "amount")

			require.EqualError(t, err, tt.expectedError)
		})
	}
}
