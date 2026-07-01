package submitswap_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap/mocks"
	swaptestdata "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap/testdata"
)

func TestUseCase_Execute_ShouldSubmitSwapWhenClientAcceptsTransaction(t *testing.T) {
	swap := chaintestdata.Swap()
	client := mocks.NewClient(t)
	client.EXPECT().SubmitSwap(t.Context(), swap).Return(nil)

	uc := submitswap.NewUseCase(client)

	output, err := uc.Execute(t.Context(), swaptestdata.Input())

	require.NoError(t, err)
	require.Equal(t, swaptestdata.Output(), output)
}

func TestUseCase_Execute_ShouldReturnErrorWhenClientSubmitFails(t *testing.T) {
	swap := chaintestdata.Swap()
	expectedErr := errors.New("transaction rejected")
	client := mocks.NewClient(t)
	client.EXPECT().SubmitSwap(t.Context(), swap).Return(expectedErr)

	uc := submitswap.NewUseCase(client)

	output, err := uc.Execute(t.Context(), swaptestdata.Input())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, submitswap.Output{}, output)
}
