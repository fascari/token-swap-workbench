//go:build integration

package quote_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/test/chainsuite"
)

type (
	QuoteSuite struct {
		chainsuite.Suite
	}
)

func TestQuoteSuite(t *testing.T) {
	suite.Run(t, new(QuoteSuite))
}

func (s *QuoteSuite) TestQuote_ShouldReturnQuote() {
	s.ChainReturnsQuote(s.ReadFile("testdata/upstream/rate.json"))

	response := s.Expect().GET("/v1/quote").
		WithQuery("in", "NEX").
		WithQuery("out", "ETH").
		WithQuery("amount", "10").
		Expect()

	response.Status(http.StatusOK)
	s.Require().JSONEq(s.ReadFile("testdata/response/quote.json"), response.Body().Raw())

	forwarded := s.LastChainRequest()
	s.Equal(http.MethodGet, forwarded.Method)
	s.Equal("/rate", forwarded.Path)
	s.Equal("NEX", forwarded.Query.Get("in"))
	s.Equal("ETH", forwarded.Query.Get("out"))
	s.Equal("10", forwarded.Query.Get("amount"))
}

func (s *QuoteSuite) TestQuote_ShouldReturnBadGatewayWhenUpstreamIsUnavailable() {
	s.ChainFailsRate(http.StatusBadGateway)

	s.Expect().GET("/v1/quote").
		WithQuery("in", "NEX").
		WithQuery("out", "ETH").
		WithQuery("amount", "10").
		Expect().
		Status(http.StatusBadGateway)
}

func (s *QuoteSuite) TestQuote_ShouldReturnBadRequestWhenUpstreamRejects() {
	s.ChainFailsRate(http.StatusUnprocessableEntity)

	s.Expect().GET("/v1/quote").
		WithQuery("in", "NEX").
		WithQuery("out", "ETH").
		WithQuery("amount", "10").
		Expect().
		Status(http.StatusBadRequest)
}
