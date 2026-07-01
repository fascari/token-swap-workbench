package chainhttp

import (
	"errors"
	"net/http"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
)

// StatusForError maps a use-case error to an HTTP status: an upstream rejection
// surfaces as 400, any other upstream failure as 502.
func StatusForError(err error) int {
	if errors.Is(err, domain.ErrUpstreamRejected) {
		return http.StatusBadRequest
	}

	return http.StatusBadGateway
}
