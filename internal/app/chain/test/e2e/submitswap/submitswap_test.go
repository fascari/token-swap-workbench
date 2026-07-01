//go:build integration

package submitswap_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/test/chainsuite"
)

type (
	SubmitSwapSuite struct {
		chainsuite.Suite
	}
)

func TestSubmitSwapSuite(t *testing.T) {
	suite.Run(t, new(SubmitSwapSuite))
}

func (s *SubmitSwapSuite) TestShouldForwardSwapAsChainTransaction() {
	s.ChainAcceptsSwap()

	response := s.Expect().POST("/v1/swaps").
		WithJSON(json.RawMessage(s.ReadFile("testdata/payload/swap.json"))).
		Expect()

	response.Status(http.StatusAccepted)
	s.Require().JSONEq(s.ReadFile("testdata/response/swap.json"), response.Body().Raw())

	forwarded := s.LastChainRequest()
	s.Equal(http.MethodPost, forwarded.Method)
	s.Equal("/transaction", forwarded.Path)
	s.Require().JSONEq(s.ReadFile("testdata/upstream/transaction.json"), forwarded.Body)
}

func (s *SubmitSwapSuite) TestShouldReturnBadRequestWhenPayloadIsInvalid() {
	s.Expect().POST("/v1/swaps").
		WithText("not-json").
		Expect().
		Status(http.StatusBadRequest)
}
