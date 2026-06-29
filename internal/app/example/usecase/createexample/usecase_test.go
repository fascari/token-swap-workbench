package createexample_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fascari/token-swap-workbench/internal/app/example/domain"
	"github.com/fascari/token-swap-workbench/internal/app/example/usecase/createexample"
	"github.com/fascari/token-swap-workbench/internal/app/example/usecase/createexample/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCase_Execute_Success(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	input := domain.Example{Name: "Test Example"}
	expected := domain.Example{ID: 1, Name: "Test Example"}

	mockRepo.EXPECT().
		Create(mock.Anything, input).
		Return(expected, nil)

	uc := createexample.New(mockRepo)
	result, err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestUseCase_Execute_Error(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	input := domain.Example{Name: "Test Example"}

	mockRepo.EXPECT().
		Create(mock.Anything, input).
		Return(domain.Example{}, errors.New("db error"))

	uc := createexample.New(mockRepo)
	result, err := uc.Execute(context.Background(), input)

	require.Error(t, err)
	require.Contains(t, err.Error(), "db error")
	require.Equal(t, domain.Example{}, result)
}
