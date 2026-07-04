package bots_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"

	handlerbots "github.com/fascari/token-swap-workbench/internal/app/chain/handler/bots"
	botsfixtures "github.com/fascari/token-swap-workbench/internal/app/chain/handler/bots/testdata"
	ucbots "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots/mocks"
	ucbotstestdata "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots/testdata"
)

const (
	botsPath = "/bots"
)

func TestHandler_Handle_ShouldCreateBots(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		uc := ucbots.New(mocks.NewClient(t))
		handler := handlerbots.New(uc)
		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, botsPath, strings.NewReader(botsfixtures.CreateRequest))
		recorder := httptest.NewRecorder()

		handler.Handle(recorder, request)
		synctest.Wait()

		require.Equal(t, http.StatusAccepted, recorder.Code)
		require.JSONEq(t, botsfixtures.CreateResponse, recorder.Body.String())

		_, err := uc.Execute(t.Context(), ucbotstestdata.StopAllInput())
		require.NoError(t, err)
		synctest.Wait()
	})
}

func TestHandler_Handle_ShouldStopAllBots(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		uc := ucbots.New(mocks.NewClient(t))
		handler := handlerbots.New(uc)
		_, err := uc.Execute(t.Context(), ucbotstestdata.CreateInput())
		require.NoError(t, err)
		synctest.Wait()

		request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, botsPath, strings.NewReader(botsfixtures.StopAllRequest))
		recorder := httptest.NewRecorder()

		handler.Handle(recorder, request)
		synctest.Wait()

		require.Equal(t, http.StatusAccepted, recorder.Code)
		require.JSONEq(t, botsfixtures.StopAllResponse, recorder.Body.String())
	})
}

func TestHandler_Handle_ShouldReturnBadRequestWhenJSONIsMalformed(t *testing.T) {
	handler := handlerbots.New(ucbots.New(mocks.NewClient(t)))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, botsPath, strings.NewReader(botsfixtures.MalformedRequest))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	requireErrorCode(t, recorder.Body.Bytes())
}

func TestHandler_Handle_ShouldReturnBadRequestWhenActionIsMissing(t *testing.T) {
	handler := handlerbots.New(ucbots.New(mocks.NewClient(t)))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, botsPath, strings.NewReader(botsfixtures.MissingActionRequest))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	requireErrorCode(t, recorder.Body.Bytes())
}

func TestHandler_Handle_ShouldReturnBadRequestWhenUseCaseValidationFails(t *testing.T) {
	handler := handlerbots.New(ucbots.New(mocks.NewClient(t)))
	request := httptest.NewRequestWithContext(t.Context(), http.MethodPost, botsPath, strings.NewReader(botsfixtures.InvalidAmountRequest))
	recorder := httptest.NewRecorder()

	handler.Handle(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	requireErrorCode(t, recorder.Body.Bytes())
}

func requireErrorCode(t *testing.T, body []byte) {
	t.Helper()

	var response struct {
		Code string `json:"code"`
	}
	require.NoError(t, json.Unmarshal(body, &response))
	require.NotEmpty(t, response.Code)
}
