package validator

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	instance = sync.OnceValue(newValidatorSet)
)

type (
	validatorSet struct {
		validate   *validator.Validate
		translator ut.Translator
	}
)

func Validate(s any) error {
	current := instance()

	err := current.validate.Struct(s)
	if err == nil {
		return nil
	}

	if fieldErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		messages := make([]string, 0, len(fieldErrors))
		for _, fieldErr := range fieldErrors {
			messages = append(messages, fieldErr.Translate(current.translator))
		}

		return errors.New(strings.Join(messages, "; "))
	}

	return err
}

func newValidatorSet() validatorSet {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}

		return name
	})

	english := en.New()
	translator, _ := ut.New(english, english).GetTranslator("en")
	_ = entranslations.RegisterDefaultTranslations(validate, translator)

	return validatorSet{validate: validate, translator: translator}
}
