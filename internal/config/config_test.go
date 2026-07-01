package config_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/fascari/token-swap-workbench/internal/config"
)

func TestLoad_ShouldReadChainBaseURLFromEnvironment(t *testing.T) {
	resetViper(t)
	t.Setenv("CHAIN_BASE_URL", "http://chain.local.test:3000")

	cfg, err := config.Load()

	require.NoError(t, err)
	require.Equal(t, "http://chain.local.test:3000", cfg.Chain.BaseURL)
}

func resetViper(t *testing.T) {
	t.Helper()

	viper.Reset()
	t.Cleanup(viper.Reset)
}
