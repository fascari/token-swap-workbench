package chainclient_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/chainclient"
	chainclienttestdata "github.com/fascari/token-swap-workbench/internal/chainclient/testdata"
	"github.com/fascari/token-swap-workbench/internal/config"
)

const (
	blockCount = 2
)

func TestClient_Quote_ShouldReturnQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/rate", r.URL.Path)
		require.Equal(t, string(chainclienttestdata.TokenUSDC), r.URL.Query().Get("in"))
		require.Equal(t, string(chainclienttestdata.TokenETH), r.URL.Query().Get("out"))
		require.Equal(t, strconv.FormatFloat(chainclienttestdata.AmountIn, 'f', -1, 64), r.URL.Query().Get("amount"))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(readFixture(t, "quote_response.json"))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	result, err := client.Quote(t.Context(), chainclienttestdata.QuoteRequest())

	require.NoError(t, err)
	require.Equal(t, chainclienttestdata.Quote(), result)
}

func TestClient_Quote_ShouldReturnUpstreamRejectedWhenChainReturns4xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad quote", http.StatusBadRequest)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	result, err := client.Quote(t.Context(), chainclienttestdata.QuoteRequest())

	require.ErrorIs(t, err, domain.ErrUpstreamRejected)
	require.Equal(t, domain.Quote{}, result)
}

func TestClient_SubmitSwap_ShouldPostSwapEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/transaction", r.URL.Path)

		body, err := os.ReadFile("testdata/swap_envelope.json")
		require.NoError(t, err)
		require.JSONEq(t, string(body), readBody(t, r))

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.SubmitSwap(t.Context(), chainclienttestdata.Swap())

	require.NoError(t, err)
}

func TestClient_SubmitSwap_ShouldReturnUpstreamRejectedWhenChainReturns4xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "rejected", http.StatusConflict)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.SubmitSwap(t.Context(), chainclienttestdata.Swap())

	require.ErrorIs(t, err, domain.ErrUpstreamRejected)
}

func TestClient_Blocks_ShouldReturnBlocks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/blocks", r.URL.Path)
		require.Equal(t, strconv.Itoa(blockCount), r.URL.Query().Get("n"))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(readFixture(t, "blocks_response.json"))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	blocks, err := client.Blocks(t.Context(), blockCount)

	require.NoError(t, err)
	require.Equal(t, chainclienttestdata.Blocks(), blocks)
}

func TestClient_Blocks_ShouldReturnUpstreamUnavailableWhenResponseIsInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"not":"a block list"}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	blocks, err := client.Blocks(t.Context(), blockCount)

	require.ErrorIs(t, err, domain.ErrUpstreamUnavailable)
	require.Empty(t, blocks)
}

func TestClient_Status_ShouldReturnUpstreamUnavailableWhenBlocksReturn5xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.Status(t.Context())

	require.ErrorIs(t, err, domain.ErrUpstreamUnavailable)
}

func newClient(t *testing.T, baseURL string) *chainclient.Client {
	t.Helper()

	client, err := chainclient.New(config.ChainConfig{BaseURL: baseURL})
	require.NoError(t, err)

	return client
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)

	return data
}

func readBody(t *testing.T, r *http.Request) string {
	t.Helper()

	data, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	return string(data)
}
