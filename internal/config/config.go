package config

import "time"

type (
	Config struct {
		App   AppConfig
		HTTP  HTTPConfig
		Log   LogConfig
		Chain ChainConfig
	}

	AppConfig struct {
		Name string
		Env  string
	}

	HTTPConfig struct {
		Port         int
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	LogConfig struct {
		Level  string
		Format string
	}

	ChainConfig struct {
		BaseURL string
	}
)
