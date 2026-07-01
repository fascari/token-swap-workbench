package domain

import "errors"

var (
	ErrUpstreamRejected    = errors.New("chain rejected the request")
	ErrUpstreamUnavailable = errors.New("chain is unavailable")
)
