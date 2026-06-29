package bootstrap

import (
	"github.com/fascari/token-swap-workbench/internal/config"
	"github.com/fascari/token-swap-workbench/pkg/logger"
)

// InitLogger initialises the global zerolog logger from the application config.
func InitLogger(cfg *config.Config) {
	logger.Init(cfg.Log.Level)
}
