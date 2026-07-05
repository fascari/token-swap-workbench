package transaction

import (
	"errors"
	"net/http"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

func errorCode(err error) int {
	if errors.Is(err, domain.ErrUpstreamRejected) {
		return http.StatusBadRequest
	}
	return http.StatusBadGateway
}
