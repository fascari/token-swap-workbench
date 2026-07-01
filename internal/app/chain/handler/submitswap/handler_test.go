package submitswap_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/handler/submitswap"
	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	submitswapuc "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap/mocks"
)

func TestHandler_Handle_ShouldSubmitSwapWhenRequestIsValid(t *testing.T) {
	body := readRequest(t)
	swap := chaintestdata.Swap()
	client := mocks.NewClient(t)
	client.EXPECT().SubmitSwap(t.Context(), swap).Return(nil)
	handler := submitswap.New(submitswapuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chain/swaps", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusAccepted, recorder.Code)
	require.JSONEq(t, string(readResponse(t)), recorder.Body.String())
}

func TestHandler_Handle_ShouldReturnBadRequestWhenPayloadIsInvalid(t *testing.T) {
	client := mocks.NewClient(t)
	handler := submitswap.New(submitswapuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chain/swaps", bytes.NewBufferString(`{`))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestHandler_Handle_ShouldReturnBadGatewayWhenSubmitFails(t *testing.T) {
	body := readRequest(t)
	swap := chaintestdata.Swap()
	client := mocks.NewClient(t)
	client.EXPECT().SubmitSwap(t.Context(), swap).Return(errors.New("transaction rejected"))
	handler := submitswap.New(submitswapuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/chain/swaps", bytes.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadGateway, recorder.Code)
}

func TestHandler_Handle_ShouldReturnBadRequestWhenPayloadFailsValidation(t *testing.T) {
	client := mocks.NewClient(t)
	handler := submitswap.New(submitswapuc.NewUseCase(client))
	request := httptest.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		"/chain/swaps",
		bytes.NewReader(readInvalidRequest(t)),
	)
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
}

func readRequest(t *testing.T) []byte {
	t.Helper()

	body, err := os.ReadFile("testdata/swap_request.json")
	require.NoError(t, err)

	return body
}

func readResponse(t *testing.T) []byte {
	t.Helper()

	body, err := os.ReadFile("testdata/response.json")
	require.NoError(t, err)

	return body
}

func readInvalidRequest(t *testing.T) []byte {
	t.Helper()

	body, err := os.ReadFile("testdata/invalid_swap_request.json")
	require.NoError(t, err)

	return body
}
