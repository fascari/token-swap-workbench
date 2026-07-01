//go:build integration

package suite

import (
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/suite"
)

// Suite is reusable, domain-agnostic end-to-end infrastructure: it starts a
// system under test from a handler, exposes an httpexpect client bound to it,
// and runs an endpoint-agnostic Upstream stub. Domain suites embed it.
type (
	Suite struct {
		suite.Suite

		Upstream *Upstream

		api *httptest.Server
	}
)

func (s *Suite) StartUpstream() {
	s.Upstream = newUpstream()
}

func (s *Suite) StartAPI(handler http.Handler) {
	s.api = httptest.NewServer(handler)
}

func (s *Suite) TearDownTest() {
	if s.api != nil {
		s.api.Close()
	}

	if s.Upstream != nil {
		s.Upstream.Close()
	}
}

func (s *Suite) Expect() *httpexpect.Expect {
	return httpexpect.Default(s.T(), s.api.URL)
}

// ReadFile reads a golden file relative to the calling test package.
func (s *Suite) ReadFile(path string) string {
	data, err := os.ReadFile(path)
	s.Require().NoError(err)

	return string(data)
}
