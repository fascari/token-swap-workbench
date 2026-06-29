package bootstrap

import (
	"github.com/fascari/token-swap-workbench/internal/config"
)

// LoadConfig delegates to the config package to load application configuration.
func LoadConfig() (*config.Config, error) {
	return config.Load()
}
