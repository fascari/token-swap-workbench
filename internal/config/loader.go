package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("APP_NAME", "token-swap-workbench")
	viper.SetDefault("APP_ENV", "development")

	viper.SetDefault("HTTP_PORT", 8080)
	viper.SetDefault("HTTP_READ_TIMEOUT", "15s")
	viper.SetDefault("HTTP_WRITE_TIMEOUT", "15s")
	viper.SetDefault("HTTP_IDLE_TIMEOUT", "60s")

	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")

	viper.SetDefault("CHAIN_BASE_URL", "http://127.0.0.1:3000")

	_ = viper.ReadInConfig()

	readTimeout, err := time.ParseDuration(viper.GetString("HTTP_READ_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT: %w", err)
	}

	writeTimeout, err := time.ParseDuration(viper.GetString("HTTP_WRITE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT: %w", err)
	}

	idleTimeout, err := time.ParseDuration(viper.GetString("HTTP_IDLE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_IDLE_TIMEOUT: %w", err)
	}

	return new(Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
		},
		HTTP: HTTPConfig{
			Port:         viper.GetInt("HTTP_PORT"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		Chain: ChainConfig{
			BaseURL: viper.GetString("CHAIN_BASE_URL"),
		},
	}), nil
}
