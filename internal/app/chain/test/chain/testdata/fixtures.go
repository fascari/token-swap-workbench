//go:build integration

package testdata

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Blocks(t *testing.T) string {
	t.Helper()

	return read(t, "blocks.json")
}

func Quote(t *testing.T) string {
	t.Helper()

	return read(t, "quote.json")
}

func Swap(t *testing.T) string {
	t.Helper()

	return read(t, "swap.json")
}

func read(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join("internal", "app", "chain", "test", "chain", "testdata", "fixtures", name)
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	return string(data)
}
