//go:build integration

package chain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/cmd/api/modules"
	"github.com/fascari/token-swap-workbench/internal/app/chain/test/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/bootstrap"
	"github.com/fascari/token-swap-workbench/internal/chainclient"
	"github.com/fascari/token-swap-workbench/internal/config"
	e2esuite "github.com/fascari/token-swap-workbench/internal/testing/integration/suite"
)

type ChainSuite struct {
	e2esuite.Suite
	upstream *upstreamAPI
}

func TestChainSuite(t *testing.T) {
	suite.Run(t, new(ChainSuite))
}

func (s *ChainSuite) SetupTest() {
	s.Suite.SetupTest()

	s.upstream = newUpstreamAPI(
		testdata.Blocks(s.T()),
		testdata.Quote(s.T()),
	)
	s.StartAPI(newWorkbenchAPI(s.T(), s.upstream.URL()))
}

func (s *ChainSuite) TearDownTest() {
	s.Suite.TearDownTest()
	s.upstream.Close()
}

func (s *ChainSuite) TestStatusShouldReturnOKWhenChainIsAvailable() {
	s.Expect().GET("/v1/chain/status").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		HasValue("status", "ok")

	req := s.upstream.LastRequest(s.T())
	require.Equal(s.T(), http.MethodGet, req.Method)
	require.Equal(s.T(), "/blocks", req.Path)
	require.Equal(s.T(), "1", req.Query.Get("n"))
}

func (s *ChainSuite) TestQuoteShouldForwardQueryAndReturnQuote() {
	s.Expect().GET("/v1/quote").
		WithQuery("in", "NEX").
		WithQuery("out", "ETH").
		WithQuery("amount", "10").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		HasValue("amount_out", 6.25)

	req := s.upstream.LastRequest(s.T())
	require.Equal(s.T(), http.MethodGet, req.Method)
	require.Equal(s.T(), "/rate", req.Path)
	require.Equal(s.T(), "NEX", req.Query.Get("in"))
	require.Equal(s.T(), "ETH", req.Query.Get("out"))
	require.Equal(s.T(), "10", req.Query.Get("amount"))
}

func (s *ChainSuite) TestSubmitSwapShouldSendChainSwapPayload() {
	s.Expect().POST("/v1/swaps").
		WithJSON(map[string]any{
			"account_id": 2,
			"in_token":   "NEX",
			"out_token":  "ETH",
			"amount_in":  10,
		}).
		Expect().
		Status(http.StatusAccepted).
		JSON().Object().
		HasValue("status", "submitted")

	req := s.upstream.LastRequest(s.T())
	require.Equal(s.T(), http.MethodPost, req.Method)
	require.Equal(s.T(), "/transaction", req.Path)
	require.JSONEq(s.T(), testdata.Swap(s.T()), req.Body)
}

func (s *ChainSuite) TestBlocksShouldForwardLimitAndReturnHistory() {
	blocks := s.Expect().GET("/v1/blocks").
		WithQuery("n", "10").
		Expect().
		Status(http.StatusOK).
		JSON().Array()

	blocks.Length().IsEqual(2)
	blocks.Value(0).Object().HasValue("id", 1)
	blocks.Value(0).Object().Value("transactions").Array().Length().IsEqual(0)
	blocks.Value(1).Object().Value("transactions").Array().Length().IsEqual(1)

	req := s.upstream.LastRequest(s.T())
	require.Equal(s.T(), http.MethodGet, req.Method)
	require.Equal(s.T(), "/blocks", req.Path)
	require.Equal(s.T(), "10", req.Query.Get("n"))
}

func (s *ChainSuite) TestQuoteShouldReturnBadGatewayWhenUpstreamFails() {
	s.upstream.FailRate(http.StatusBadGateway)

	s.Expect().GET("/v1/quote").
		WithQuery("in", "NEX").
		WithQuery("out", "ETH").
		WithQuery("amount", "10").
		Expect().
		Status(http.StatusBadGateway)
}

func (s *ChainSuite) TestBlocksShouldReturnBadGatewayWhenUpstreamPayloadIsMalformed() {
	s.upstream.MalformedBlocks()

	s.Expect().GET("/v1/blocks").
		WithQuery("n", "10").
		Expect().
		Status(http.StatusBadGateway)
}

func (s *ChainSuite) TestFrontendShouldServeIndex() {
	s.Expect().GET("/").
		Expect().
		Status(http.StatusOK).
		ContentType("text/html")
}

func newWorkbenchAPI(t *testing.T, chainBaseURL string) http.Handler {
	t.Helper()

	client, err := chainclient.New(config.ChainConfig{
		BaseURL: chainBaseURL,
	})
	require.NoError(t, err)

	return bootstrap.NewRouter(modules.NewChainModule(client))
}
