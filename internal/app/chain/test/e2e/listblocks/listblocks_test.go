//go:build integration

package listblocks_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/test/chainsuite"
)

type (
	ListBlocksSuite struct {
		chainsuite.Suite
	}
)

func TestListBlocksSuite(t *testing.T) {
	suite.Run(t, new(ListBlocksSuite))
}

func (s *ListBlocksSuite) TestShouldReturnHistoryWhenForwardingLimit() {
	s.ChainReturnsBlocks(s.ReadFile("testdata/upstream/blocks.json"))

	response := s.Expect().GET("/v1/blocks").
		WithQuery("n", "10").
		Expect()

	response.Status(http.StatusOK)
	s.Require().JSONEq(s.ReadFile("testdata/response/blocks.json"), response.Body().Raw())

	forwarded := s.LastChainRequest()
	s.Equal(http.MethodGet, forwarded.Method)
	s.Equal("/blocks", forwarded.Path)
	s.Equal("10", forwarded.Query.Get("n"))
}

func (s *ListBlocksSuite) TestShouldReturnBadGatewayWhenUpstreamPayloadIsMalformed() {
	s.ChainReturnsMalformedBlocks()

	s.Expect().GET("/v1/blocks").
		WithQuery("n", "10").
		Expect().
		Status(http.StatusBadGateway)
}
