package status_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/status"
	statusuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status/mocks"
)

func TestHandler_Handle_ShouldReturnStatusWhenChainIsAvailable(t *testing.T) {
	client := mocks.NewClient(t)
	client.EXPECT().Status(t.Context()).Return(nil)
	handler := status.New(statusuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/chain/status", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, string(readResponse(t)), recorder.Body.String())
}

func TestHandler_Handle_ShouldReturnBadGatewayWhenChainStatusFails(t *testing.T) {
	client := mocks.NewClient(t)
	client.EXPECT().Status(t.Context()).Return(errors.New("chain unavailable"))
	handler := status.New(statusuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/chain/status", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadGateway, recorder.Code)
}

func readResponse(t *testing.T) []byte {
	t.Helper()

	body, err := os.ReadFile("testdata/response.json")
	require.NoError(t, err)

	return body
}
