//go:build integration

package suite

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	APIURL string

	api *httptest.Server
}

func (s *Suite) SetupTest() {
	s.T().Chdir(repositoryRoot(s.T()))
}

func (s *Suite) TearDownTest() {
	if s.api != nil {
		s.api.Close()
	}
}

func (s *Suite) StartAPI(handler http.Handler) {
	s.api = httptest.NewServer(handler)
	s.APIURL = s.api.URL
}

func (s *Suite) Expect() *httpexpect.Expect {
	return httpexpect.Default(s.T(), s.APIURL)
}

func repositoryRoot(t require.TestingT) string {
	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		require.NotEqual(t, dir, parent)

		dir = parent
	}
}
