//go:build integration

package suite

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// Upstream is an endpoint-agnostic stand-in for an external HTTP service. It
// records every request and serves responses registered per method and path.
type (
	Upstream struct {
		server *httptest.Server

		mu       sync.Mutex
		requests []Request
		stubs    map[string]stub
	}

	Request struct {
		Method string
		Path   string
		Query  url.Values
		Body   string
	}

	stub struct {
		status  int
		body    string
		handler http.HandlerFunc
	}
)

func newUpstream() *Upstream {
	upstream := &Upstream{stubs: make(map[string]stub)}
	upstream.server = httptest.NewServer(http.HandlerFunc(upstream.serve))

	return upstream
}

func (u *Upstream) URL() string {
	return u.server.URL
}

func (u *Upstream) Close() {
	u.server.Close()
}

// Stub registers a canned JSON response for a method and path.
func (u *Upstream) Stub(method, path string, status int, body string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.stubs[stubKey(method, path)] = stub{status: status, body: body}
}

// StubHandler registers a custom handler for responses a canned body cannot express.
func (u *Upstream) StubHandler(method, path string, handler http.HandlerFunc) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.stubs[stubKey(method, path)] = stub{handler: handler}
}

func (u *Upstream) LastRequest(t *testing.T) Request {
	t.Helper()

	u.mu.Lock()
	defer u.mu.Unlock()

	require.NotEmpty(t, u.requests, "expected the system under test to call the upstream")

	return u.requests[len(u.requests)-1]
}

func (u *Upstream) serve(w http.ResponseWriter, r *http.Request) {
	u.capture(r)

	u.mu.Lock()
	registered, ok := u.stubs[stubKey(r.Method, r.URL.Path)]
	u.mu.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if registered.handler != nil {
		registered.handler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(registered.status)
	_, _ = io.WriteString(w, registered.body)
}

func (u *Upstream) capture(r *http.Request) {
	body := ""
	if r.Body != nil {
		if data, err := io.ReadAll(r.Body); err == nil {
			body = string(data)
		}
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	u.requests = append(u.requests, Request{
		Method: r.Method,
		Path:   r.URL.Path,
		Query:  r.URL.Query(),
		Body:   body,
	})
}

func stubKey(method, path string) string {
	return method + " " + path
}
