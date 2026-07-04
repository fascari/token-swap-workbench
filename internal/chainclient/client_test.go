package chainclient_test

import (
	"io"
	"net/http"
	"net/http/httptest"
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
		_, err := w.Write([]byte(chainclienttestdata.QuoteResponse))
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

func TestClient_SubmitTransaction_ShouldPostSendEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/transaction", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.JSONEq(t, chainclienttestdata.SendEnvelope, readBody(t, r))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	client, err := chainclient.New(config.ChainConfig{BaseURL: server.URL})
	require.NoError(t, err)

	err = client.SubmitTransaction(t.Context(), chainclienttestdata.SendSubmission())

	require.NoError(t, err)
}

func TestClient_SubmitTransaction_ShouldPostSwapEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/transaction", r.URL.Path)

		require.JSONEq(t, chainclienttestdata.SwapEnvelope, readBody(t, r))

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.SubmitTransaction(t.Context(), chainclienttestdata.SwapSubmission())

	require.NoError(t, err)
}

func TestClient_SubmitTransaction_ShouldReturnUpstreamRejectedWhenChainReturns4xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "rejected", http.StatusConflict)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.SubmitTransaction(t.Context(), chainclienttestdata.SwapSubmission())

	require.ErrorIs(t, err, domain.ErrUpstreamRejected)
}

func TestClient_Blocks_ShouldReturnBlocks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/blocks", r.URL.Path)
		require.Equal(t, strconv.Itoa(blockCount), r.URL.Query().Get("n"))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(chainclienttestdata.BlocksResponse))
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

func readBody(t *testing.T, r *http.Request) string {
	t.Helper()

	data, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	return string(data)
}
