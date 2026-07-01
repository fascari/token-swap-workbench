package submitswap_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/submitswap/mocks"
)

func TestUseCase_Execute_ShouldSubmitSwapWhenClientAcceptsTransaction(t *testing.T) {
	swap := chaintestdata.Swap()
	client := mocks.NewClient(t)
	client.EXPECT().SubmitSwap(t.Context(), swap).Return(nil)

	uc := submitswap.NewUseCase(client)

	output, err := uc.Execute(t.Context(), submitswap.Input{
		AccountID: swap.AccountID,
		InToken:   swap.InToken,
		OutToken:  swap.OutToken,
		AmountIn:  swap.AmountIn,
	})

	require.NoError(t, err)
	require.Equal(t, submitswap.Output{Status: "submitted"}, output)
}

func TestUseCase_Execute_ShouldReturnErrorWhenClientSubmitFails(t *testing.T) {
	swap := chaintestdata.Swap()
	expectedErr := errors.New("transaction rejected")
	client := mocks.NewClient(t)
	client.EXPECT().SubmitSwap(t.Context(), swap).Return(expectedErr)

	uc := submitswap.NewUseCase(client)

	output, err := uc.Execute(t.Context(), submitswap.Input{
		AccountID: swap.AccountID,
		InToken:   swap.InToken,
		OutToken:  swap.OutToken,
		AmountIn:  swap.AmountIn,
	})

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, submitswap.Output{}, output)
}
