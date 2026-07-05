package bots

import (
	"errors"
	"net/http"

	botsuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots"
)

func errorCode(err error) int {
	if errors.Is(err, botsuc.ErrInvalidAction) || errors.Is(err, botsuc.ErrInvalidAmount) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
