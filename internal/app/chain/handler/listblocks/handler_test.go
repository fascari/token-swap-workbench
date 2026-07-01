package listblocks_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/listblocks"
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	listblocksuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks/mocks"
)

const (
	blockCount = 2
)

func TestHandler_Handle_ShouldReturnBlocksWhenRequestIsValid(t *testing.T) {
	expectedBlocks := chaintestdata.Blocks()
	client := mocks.NewClient(t)
	client.EXPECT().Blocks(t.Context(), blockCount).Return(expectedBlocks, nil)
	handler := listblocks.New(listblocksuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/chain/blocks?n=2", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.JSONEq(t, string(readResponse(t)), recorder.Body.String())
}

func TestHandler_Handle_ShouldReturnBadRequestWhenCountIsInvalid(t *testing.T) {
	client := mocks.NewClient(t)
	handler := listblocks.New(listblocksuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/chain/blocks?n=bad", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Handle_ShouldReturnBadGatewayWhenClientListFails(t *testing.T) {
	client := mocks.NewClient(t)
	client.EXPECT().Blocks(t.Context(), blockCount).Return(nil, errors.New("blocks unavailable"))
	handler := listblocks.New(listblocksuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/chain/blocks?n=2", nil)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadGateway, recorder.Code)
}

func readResponse(t *testing.T) []byte {
	t.Helper()

	body, err := os.ReadFile("testdata/blocks_response.json")
	require.NoError(t, err)

	return body
}
