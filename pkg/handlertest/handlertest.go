package handlertest

import (
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/suite"
)

// Suite is reusable, domain-agnostic infrastructure for unit-testing HTTP
// handlers in isolation. It serves a handler with an in-memory recorder and
// captures the response so domain handler suites can assert on it. Fixtures are
// provided by the caller (embedded via go:embed in the domain testdata package).
type (
	Suite struct {
		suite.Suite

		ResponseCode int
		ResponseBody string
	}
)

func (s *Suite) Serve(handler http.HandlerFunc, method string, target string, opts ...OptionFunc) {
	s.T().Helper()

	request := s.buildRequest(method, target, opts...)
	recorder := httptest.NewRecorder()

	handler(recorder, request)

	s.ResponseCode = recorder.Code
	s.ResponseBody = recorder.Body.String()
}

func (s *Suite) RequireStatus(status int) {
	s.T().Helper()

	s.Require().Equal(status, s.ResponseCode)
}

func (s *Suite) RequireJSONResponse(status int, expected string) {
	s.T().Helper()

	s.Require().Equal(status, s.ResponseCode)
	s.Require().JSONEq(expected, s.ResponseBody)
}
