package bootstrap

import (
	"github.com/fascari/token-swap-workbench/internal/config"
	"github.com/fascari/token-swap-workbench/pkg/logger"
)

func InitLogger(cfg *config.Config) {
	logger.Init(cfg.Log.Level)
}
