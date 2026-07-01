//go:build integration

package status_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/test/chainsuite"
)

type (
	StatusSuite struct {
		chainsuite.Suite
	}
)

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}

func (s *StatusSuite) TestShouldReturnOKWhenChainIsAvailable() {
	s.ChainReturnsBlocks(s.ReadFile("testdata/upstream/blocks.json"))

	response := s.Expect().GET("/v1/chain/status").Expect()

	response.Status(http.StatusOK)
	s.Require().JSONEq(s.ReadFile("testdata/response/status.json"), response.Body().Raw())

	forwarded := s.LastChainRequest()
	s.Equal(http.MethodGet, forwarded.Method)
	s.Equal("/blocks", forwarded.Path)
	s.Equal("1", forwarded.Query.Get("n"))
}

func (s *StatusSuite) TestShouldReturnBadGatewayWhenChainIsUnavailable() {
	s.ChainReturnsMalformedBlocks()

	s.Expect().GET("/v1/chain/status").
		Expect().
		Status(http.StatusBadGateway)
}
