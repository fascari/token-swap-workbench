//go:build integration

package chainsuite

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/cmd/api/modules"
	"github.com/fascari/token-swap-workbench/internal/bootstrap"
	"github.com/fascari/token-swap-workbench/internal/chainclient"
	"github.com/fascari/token-swap-workbench/internal/config"
	e2e "github.com/fascari/token-swap-workbench/internal/testing/integration/suite"
)

// Suite embeds the generic e2e infrastructure and wires the workbench router
// to the upstream stub, exposing chain helpers in domain terms.
type (
	Suite struct {
		e2e.Suite
	}
)

func (s *Suite) SetupTest() {
	s.StartUpstream()

	client, err := chainclient.New(config.ChainConfig{BaseURL: s.Upstream.URL()})
	require.NoError(s.T(), err)

	s.StartAPI(bootstrap.NewRouter(modules.NewChainModule(client)))
}

func (s *Suite) ChainReturnsQuote(body string) {
	s.Upstream.Stub(http.MethodGet, "/rate", http.StatusOK, body)
}

func (s *Suite) ChainReturnsBlocks(body string) {
	s.Upstream.Stub(http.MethodGet, "/blocks", http.StatusOK, body)
}

func (s *Suite) ChainAcceptsSwap() {
	s.Upstream.Stub(http.MethodPost, "/transaction", http.StatusOK, "")
}

func (s *Suite) ChainFailsRate(status int) {
	s.Upstream.Stub(http.MethodGet, "/rate", status, `{"error":"rate unavailable"}`)
}

func (s *Suite) ChainReturnsMalformedBlocks() {
	s.Upstream.StubHandler(http.MethodGet, "/blocks", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{"))
	})
}

func (s *Suite) LastChainRequest() e2e.Request {
	return s.Upstream.LastRequest(s.T())
}
