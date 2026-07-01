//go:build integration

package chain

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	upstreamAPI struct {
		server   *httptest.Server
		requests []upstreamRequest

		blocks     string
		quote      string
		rateStatus int
		badBlocks  bool
		mu         sync.Mutex
	}

	upstreamRequest struct {
		Method string
		Path   string
		Query  url.Values
		Body   string
	}
)

func newUpstreamAPI(blocks string, quote string) *upstreamAPI {
	api := &upstreamAPI{
		blocks:     blocks,
		quote:      quote,
		rateStatus: http.StatusOK,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/blocks", api.handleBlocks)
	mux.HandleFunc("/rate", api.handleRate)
	mux.HandleFunc("/transaction", api.handleTransaction)

	api.server = httptest.NewServer(mux)

	return api
}

func (a *upstreamAPI) Close() {
	a.server.Close()
}

func (a *upstreamAPI) URL() string {
	return a.server.URL
}

func (a *upstreamAPI) LastRequest(t *testing.T) upstreamRequest {
	t.Helper()

	a.mu.Lock()
	defer a.mu.Unlock()

	require.NotEmpty(t, a.requests)

	return a.requests[len(a.requests)-1]
}

func (a *upstreamAPI) FailRate(status int) {
	a.rateStatus = status
}

func (a *upstreamAPI) MalformedBlocks() {
	a.badBlocks = true
}

func (a *upstreamAPI) handleBlocks(w http.ResponseWriter, r *http.Request) {
	a.capture(r)

	if a.badBlocks {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{"))
		return
	}

	writeJSON(w, http.StatusOK, a.blocks)
}

func (a *upstreamAPI) handleRate(w http.ResponseWriter, r *http.Request) {
	a.capture(r)

	if a.rateStatus != http.StatusOK {
		writeJSON(w, a.rateStatus, map[string]string{"error": "rate unavailable"})
		return
	}

	writeJSON(w, http.StatusOK, a.quote)
}

func (a *upstreamAPI) handleTransaction(w http.ResponseWriter, r *http.Request) {
	a.capture(r)
	w.WriteHeader(http.StatusOK)
}

func (a *upstreamAPI) capture(r *http.Request) {
	body := ""
	if r.Body != nil {
		data, err := io.ReadAll(r.Body)
		if err == nil {
			body = string(data)
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.requests = append(a.requests, upstreamRequest{
		Method: r.Method,
		Path:   r.URL.Path,
		Query:  r.URL.Query(),
		Body:   body,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if value, ok := data.(string); ok {
		_, _ = w.Write([]byte(value))
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}
