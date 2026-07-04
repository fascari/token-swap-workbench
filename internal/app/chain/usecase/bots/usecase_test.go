package bots_test

import (
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots/mocks"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/bots/testdata"
)

func TestUseCase_Execute_ShouldCreateBots(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		uc := bots.New(mocks.NewClient(t))

		output, err := uc.Execute(t.Context(), testdata.CreateInput())
		synctest.Wait()

		require.NoError(t, err)
		require.Equal(t, testdata.CreatedOutput(), output)

		_, err = uc.Execute(t.Context(), testdata.StopAllInput())
		require.NoError(t, err)
		synctest.Wait()
	})
}

func TestUseCase_Execute_ShouldStopRequestedBots(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		uc := bots.New(mocks.NewClient(t))
		_, err := uc.Execute(t.Context(), testdata.CreateInput())
		require.NoError(t, err)
		synctest.Wait()

		output, err := uc.Execute(t.Context(), testdata.StopInput())
		synctest.Wait()

		require.NoError(t, err)
		require.Equal(t, testdata.StoppedOutput(), output)
	})
}

func TestUseCase_Execute_ShouldStopAllBots(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		uc := bots.New(mocks.NewClient(t))
		_, err := uc.Execute(t.Context(), testdata.CreateInput())
		require.NoError(t, err)
		synctest.Wait()

		output, err := uc.Execute(t.Context(), testdata.StopAllInput())
		synctest.Wait()

		require.NoError(t, err)
		require.Equal(t, testdata.StoppedAllOutput(), output)
	})
}

func TestUseCase_Execute_ShouldRejectInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		input   bots.Input
		wantErr error
	}{
		{
			name:    "should reject invalid action",
			input:   testdata.InvalidActionInput(),
			wantErr: bots.ErrInvalidAction,
		},
		{
			name:    "should reject all with create action",
			input:   testdata.AllWithCreateInput(),
			wantErr: bots.ErrInvalidAction,
		},
		{
			name:    "should reject zero amount",
			input:   testdata.ZeroAmountInput(),
			wantErr: bots.ErrInvalidAmount,
		},
		{
			name:    "should reject amount above maximum",
			input:   testdata.TooLargeAmountInput(),
			wantErr: bots.ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := bots.New(mocks.NewClient(t))

			_, err := uc.Execute(t.Context(), tt.input)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
