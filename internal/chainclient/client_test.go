package chainclient_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/chainclient"
	"github.com/fascari/token-swap-workbench/internal/config"
)

const (
	quoteAmountIn       = 12.5
	expectedQuoteAmount = 6.25
	blockCount          = 2
)

func TestClient_Quote_ShouldReturnAmountOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/rate", r.URL.Path)
		require.Equal(t, string(chainclient.TokenUSDC), r.URL.Query().Get("in"))
		require.Equal(t, string(chainclient.TokenETH), r.URL.Query().Get("out"))
		require.Equal(t, "12.5", r.URL.Query().Get("amount"))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(readFixture(t, "quote_response.json"))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	result, err := client.Quote(t.Context(), chainclient.QuoteRequest{
		InToken:  chainclient.TokenUSDC,
		OutToken: chainclient.TokenETH,
		Amount:   quoteAmountIn,
	})

	require.NoError(t, err)
	require.Equal(t, expectedQuoteAmount, result.AmountOut)
}

func TestClient_Quote_ShouldReturnErrorWhenChainRejectsRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "bad quote", http.StatusBadRequest)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	result, err := client.Quote(t.Context(), chainclient.QuoteRequest{
		InToken:  chainclient.TokenUSDC,
		OutToken: chainclient.TokenETH,
		Amount:   quoteAmountIn,
	})

	require.Error(t, err)
	require.Equal(t, chainclient.QuoteResponse{}, result)
}

func TestClient_SubmitSwap_ShouldPostSwapTransaction(t *testing.T) {
	expectedSwap := decodeFixture[chainclient.SwapRequest](t, "swap_request.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/transaction", r.URL.Path)

		var body struct {
			Swap chainclient.SwapRequest `json:"Swap"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, expectedSwap, body.Swap)

		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.SubmitSwap(t.Context(), expectedSwap)

	require.NoError(t, err)
}

func TestClient_SubmitSwap_ShouldReturnErrorWhenChainRejectsTransaction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "rejected", http.StatusConflict)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.SubmitSwap(t.Context(), decodeFixture[chainclient.SwapRequest](t, "swap_request.json"))

	require.Error(t, err)
}

func TestClient_Blocks_ShouldReturnBlocks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/blocks", r.URL.Path)
		require.Equal(t, "2", r.URL.Query().Get("n"))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(readFixture(t, "blocks_response.json"))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	blocks, err := client.Blocks(t.Context(), blockCount)

	require.NoError(t, err)
	require.Equal(t, decodeFixture[[]chainclient.Block](t, "blocks_response.json"), blocks)
}

func TestClient_Blocks_ShouldReturnErrorWhenResponseIsInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"not":"a block list"}`))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	blocks, err := client.Blocks(t.Context(), blockCount)

	require.Error(t, err)
	require.Empty(t, blocks)
}

func TestClient_Status_ShouldReturnErrorWhenBlocksUnavailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	t.Cleanup(server.Close)

	client := newClient(t, server.URL)

	err := client.Status(t.Context())

	require.Error(t, err)
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

func decodeFixture[T any](t *testing.T, name string) T {
	t.Helper()

	var value T
	require.NoError(t, json.Unmarshal(readFixture(t, name), &value))

	return value
}
