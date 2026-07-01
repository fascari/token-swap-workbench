package status_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status/mocks"
	statustestdata "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/status/testdata"
)

func TestUseCase_Execute_ShouldReturnOKWhenClientStatusSucceeds(t *testing.T) {
	client := mocks.NewClient(t)
	client.EXPECT().Status(t.Context()).Return(nil)

	uc := status.NewUseCase(client)

	output, err := uc.Execute(t.Context(), status.Input{})

	require.NoError(t, err)
	require.Equal(t, statustestdata.Output(), output)
}

func TestUseCase_Execute_ShouldReturnErrorWhenClientStatusFails(t *testing.T) {
	expectedErr := errors.New("chain unavailable")
	client := mocks.NewClient(t)
	client.EXPECT().Status(t.Context()).Return(expectedErr)

	uc := status.NewUseCase(client)

	output, err := uc.Execute(t.Context(), status.Input{})

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, status.Output{}, output)
}
