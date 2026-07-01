package httpparam

import (
	"errors"
	"fmt"
	"strconv"
)

func Required(value, field string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("%s is required", field)
	}

	return value, nil
}

func PositiveInt(value, field string) (int, error) {
	number, err := parse(value, field, strconv.Atoi, "an integer")
	if err != nil {
		return 0, err
	}

	if number <= 0 {
		return 0, errors.New(field + " must be greater than 0")
	}

	return number, nil
}

func PositiveFloat(value, field string) (float64, error) {
	number, err := parse(value, field, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	}, "a number")
	if err != nil {
		return 0, err
	}

	if number <= 0 {
		return 0, errors.New(field + " must be greater than 0")
	}

	return number, nil
}

func parse[T any](value, field string, convert func(string) (T, error), kind string) (T, error) {
	var zero T

	if value == "" {
		return zero, fmt.Errorf("%s is required", field)
	}

	parsed, err := convert(value)
	if err != nil {
		return zero, fmt.Errorf("%s must be %s: %w", field, kind, err)
	}

	return parsed, nil
}
