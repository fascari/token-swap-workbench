package transaction_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/app/chain/domain"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction"
	"github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction/mocks"
	transactiontestdata "github.com/fascari/token-swap-workbench/internal/app/chain/usecase/transaction/testdata"
)

func TestUseCase_Execute_ShouldSubmitTransaction(t *testing.T) {
	tests := []struct {
		name       string
		submission domain.TransactionSubmission
	}{
		{
			name:       "should submit explicit send transaction",
			submission: transactiontestdata.SendSubmission(),
		},
		{
			name:       "should submit explicit swap transaction",
			submission: transactiontestdata.SwapSubmission(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mocks.NewClient(t)
			client.EXPECT().SubmitTransaction(t.Context(), tt.submission).Return(nil)
			uc := transaction.New(client)

			output, err := uc.Execute(t.Context(), tt.submission)

			require.NoError(t, err)
			require.Equal(t, transactiontestdata.Output(), output)
		})
	}
}

func TestUseCase_Execute_ShouldReturnErrorWhenClientSubmitFails(t *testing.T) {
	expectedErr := errors.New("transaction rejected")
	client := mocks.NewClient(t)
	client.EXPECT().SubmitTransaction(t.Context(), transactiontestdata.SwapSubmission()).Return(expectedErr)
	uc := transaction.New(client)

	output, err := uc.Execute(t.Context(), transactiontestdata.SwapSubmission())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, transaction.Output{}, output)
}
