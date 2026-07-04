package transaction_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/transaction"
	transactionfixtures "github.com/fascari/token-swap-workbench/internal/app/chain/handler/transaction/testdata"
	transactionuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction/mocks"
	transactiontestdata "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction/testdata"
	"github.com/fascari/token-swap-workbench/pkg/handlertest"
)

type (
	TransactionSuite struct {
		handlertest.Suite
	}
)

func TestTransactionSuite(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

func (s *TransactionSuite) TestHandler_Handle_ShouldSubmitTransaction() {
	client := mocks.NewClient(s.T())
	client.EXPECT().SubmitTransaction(s.T().Context(), transactiontestdata.SwapSubmission()).Return(nil)
	handler := transaction.New(transactionuc.New(client))

	s.Serve(handler.Handle, http.MethodPost, "/chain/transactions", handlertest.WithBody(transactionfixtures.Request))

	s.RequireJSONResponse(http.StatusAccepted, transactionfixtures.Response)
}

func (s *TransactionSuite) TestHandler_Handle_ShouldReturnBadRequestWhenPayloadIsMalformed() {
	handler := transaction.New(transactionuc.New(mocks.NewClient(s.T())))

	s.Serve(handler.Handle, http.MethodPost, "/chain/transactions", handlertest.WithBody("{"))

	s.RequireStatus(http.StatusBadRequest)
}

func (s *TransactionSuite) TestHandler_Handle_ShouldReturnBadGatewayWhenSubmitFails() {
	client := mocks.NewClient(s.T())
	client.EXPECT().SubmitTransaction(s.T().Context(), transactiontestdata.SwapSubmission()).Return(errors.New("transaction rejected"))
	handler := transaction.New(transactionuc.New(client))

	s.Serve(handler.Handle, http.MethodPost, "/chain/transactions", handlertest.WithBody(transactionfixtures.Request))

	s.RequireStatus(http.StatusBadGateway)
}

func (s *TransactionSuite) TestHandler_Handle_ShouldReturnBadRequestWhenPayloadFailsValidation() {
	handler := transaction.New(transactionuc.New(mocks.NewClient(s.T())))

	s.Serve(handler.Handle, http.MethodPost, "/chain/transactions", handlertest.WithBody(transactionfixtures.InvalidRequest))

	s.RequireStatus(http.StatusBadRequest)
}
