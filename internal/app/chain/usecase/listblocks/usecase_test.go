package listblocks_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	chaintestdata "github.com/fascari/token-swap-workbench/internal/app/chain/testdata"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/listblocks/mocks"
)

const blockCount = 2

func TestUseCase_Execute_ShouldReturnBlocksWhenClientListsBlocks(t *testing.T) {
	expectedBlocks := chaintestdata.Blocks()
	client := mocks.NewClient(t)
	client.EXPECT().Blocks(t.Context(), blockCount).Return(expectedBlocks, nil)

	uc := listblocks.NewUseCase(client)

	output, err := uc.Execute(t.Context(), listblocks.Input{Count: blockCount})

	require.NoError(t, err)
	require.Equal(t, listblocks.Output{Blocks: expectedBlocks}, output)
}

func TestUseCase_Execute_ShouldReturnErrorWhenClientListFails(t *testing.T) {
	expectedErr := errors.New("blocks unavailable")
	client := mocks.NewClient(t)
	client.EXPECT().Blocks(t.Context(), blockCount).Return(nil, expectedErr)

	uc := listblocks.NewUseCase(client)

	output, err := uc.Execute(t.Context(), listblocks.Input{Count: blockCount})

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, listblocks.Output{}, output)
}
