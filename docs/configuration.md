# Configuration

Configuration uses [Viper](https://github.com/spf13/viper) with multiple sources.

## Precedence

1. **Environment variables** (highest priority)
2. **config.yaml** file
3. **Default values** in code

## Environment Variables

All config keys can be set via environment variables using the `APP_` prefix:

```bash
APP_HTTP_PORT=3000
APP_LOG_LEVEL=debug
```

Viper converts uppercase with underscores to nested keys:
- `APP_HTTP_PORT` → `http.port`
- `APP_LOG_LEVEL` → `log.level`

## .env File

Create a `.env` file in the project root for local development:

```bash
cp .env.example .env
```

The mise task runner automatically loads `.env` before running commands.

## config.yaml

Optional YAML file for structured config:

```yaml
app:
  name: token-swap-workbench
  env: development

http:
  port: 8080

log:
  level: info
  format: json
```

## Code Example

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    App      AppConfig
    HTTP     HTTPConfig
    Log      LogConfig
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()
    viper.SetEnvPrefix("APP")

    viper.SetDefault("http.port", 8080)
    viper.SetDefault("log.level", "info")

    _ = viper.ReadInConfig() // optional

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```
