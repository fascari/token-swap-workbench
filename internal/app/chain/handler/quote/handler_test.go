package quote_test

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/quote"
	quotefixtures "github.com/fascari/token-swap-workbench/internal/app/chain/handler/quote/testdata"
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	quoteuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/quote/mocks"
	"github.com/fascari/token-swap-workbench/pkg/handlertest"
)

type (
	QuoteSuite struct {
		handlertest.Suite
	}
)

func TestQuoteSuite(t *testing.T) {
	suite.Run(t, new(QuoteSuite))
}

func (s *QuoteSuite) TestHandler_Handle_ShouldReturnQuote() {
	quoteRequest := chaintestdata.QuoteRequest()
	client := mocks.NewClient(s.T())
	client.EXPECT().Quote(s.T().Context(), quoteRequest).Return(chaintestdata.Quote(), nil)
	handler := quote.New(quoteuc.New(client))

	s.Serve(handler.Handle, http.MethodGet, "/chain/quote",
		handlertest.WithQuery("in", string(quoteRequest.InToken)),
		handlertest.WithQuery("out", string(quoteRequest.OutToken)),
		handlertest.WithQuery("amount", strconv.FormatFloat(quoteRequest.Amount, 'f', -1, 64)),
	)

	s.RequireJSONResponse(http.StatusOK, quotefixtures.Response)
}

func (s *QuoteSuite) TestHandler_Handle_ShouldReturnBadRequestWhenAmountIsMissing() {
	handler := quote.New(quoteuc.New(mocks.NewClient(s.T())))

	s.Serve(handler.Handle, http.MethodGet, "/chain/quote",
		handlertest.WithQuery("in", string(chaintestdata.TokenUSDC)),
		handlertest.WithQuery("out", string(chaintestdata.TokenETH)),
	)

	s.RequireStatus(http.StatusBadRequest)
}

func (s *QuoteSuite) TestHandler_Handle_ShouldReturnBadGatewayWhenQuoteFails() {
	quoteRequest := chaintestdata.QuoteRequest()
	client := mocks.NewClient(s.T())
	client.EXPECT().Quote(s.T().Context(), quoteRequest).Return(chaintestdata.Quote(), errors.New("quote unavailable"))
	handler := quote.New(quoteuc.New(client))

	s.Serve(handler.Handle, http.MethodGet, "/chain/quote",
		handlertest.WithQuery("in", string(quoteRequest.InToken)),
		handlertest.WithQuery("out", string(quoteRequest.OutToken)),
		handlertest.WithQuery("amount", strconv.FormatFloat(quoteRequest.Amount, 'f', -1, 64)),
	)

	s.RequireStatus(http.StatusBadGateway)
}
