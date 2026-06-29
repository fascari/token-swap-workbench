package createexample

import (
	"github.com/fascari/token-swap-workbench/internal/app/example/domain"
	"github.com/fascari/token-swap-workbench/pkg/validator"
)

type (
	InputPayload struct {
		Name string `json:"name" validate:"required,min=1,max=255"`
	}

	OutputPayload struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
)

func (p InputPayload) Validate() error {
	return validator.Validate(p)
}

func (p InputPayload) ToDomain() domain.Example {
	return domain.Example{Name: p.Name}
}

func ToOutputPayload(e domain.Example) OutputPayload {
	return OutputPayload{
		ID:   e.ID,
		Name: e.Name,
	}
}
