package status_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/status"
	statusfixtures "github.com/fascari/token-swap-workbench/internal/app/chain/handler/status/testdata"
	statusuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status/mocks"
	"github.com/fascari/token-swap-workbench/pkg/handlertest"
)

type (
	StatusSuite struct {
		handlertest.Suite
	}
)

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}

func (s *StatusSuite) TestHandler_Handle_ShouldReturnStatusWhenChainIsAvailable() {
	client := mocks.NewClient(s.T())
	client.EXPECT().Status(s.T().Context()).Return(nil)
	handler := status.New(statusuc.New(client))

	s.Serve(handler.Handle, http.MethodGet, "/chain/status")

	s.RequireJSONResponse(http.StatusOK, statusfixtures.Response)
}

func (s *StatusSuite) TestHandler_Handle_ShouldReturnBadGatewayWhenChainStatusFails() {
	client := mocks.NewClient(s.T())
	client.EXPECT().Status(s.T().Context()).Return(errors.New("chain unavailable"))
	handler := status.New(statusuc.New(client))

	s.Serve(handler.Handle, http.MethodGet, "/chain/status")

	s.RequireStatus(http.StatusBadGateway)
}
