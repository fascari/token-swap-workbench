package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	instance *validator.Validate
	once     sync.Once
)

func Validator() *validator.Validate {
	once.Do(func() {
		instance = validator.New()
	})
	return instance
}

func Validate(s any) error {
	return Validator().Struct(s)
}
