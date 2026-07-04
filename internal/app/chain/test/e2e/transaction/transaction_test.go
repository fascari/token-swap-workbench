//go:build integration

package transaction_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/test/chainsuite"
)

type (
	TransactionSuite struct {
		chainsuite.Suite
	}
)

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) TestSubmitTransaction_ShouldForwardSwapAsChainTransaction() {
	s.ChainAcceptsSwap()

	response := s.Expect().POST("/v1/transactions").
		WithJSON(json.RawMessage(s.ReadFile("testdata/payload/transaction.json"))).
		Expect()

	response.Status(http.StatusAccepted)
	s.Require().JSONEq(s.ReadFile("testdata/response/transaction.json"), response.Body().Raw())

	forwarded := s.LastChainRequest()
	s.Equal(http.MethodPost, forwarded.Method)
	s.Equal("/transaction", forwarded.Path)
	s.Require().JSONEq(s.ReadFile("testdata/upstream/transaction.json"), forwarded.Body)
}

func (s *TransactionSuite) TestSubmitTransaction_ShouldReturnBadRequestWhenPayloadIsInvalid() {
	s.Expect().POST("/v1/transactions").
		WithText("not-json").
		Expect().
		Status(http.StatusBadRequest)
}
