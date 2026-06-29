package config

import "time"

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	Log  LogConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type LogConfig struct {
	Level  string
	Format string
}
