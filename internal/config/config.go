package config

import (
	"fmt"
	"time"

	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	koanfenv "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config represents the application configuration
type Config struct {
	Server       ServerConfig     `koanf:"server"`
	Log          LogConfig        `koanf:"log"`
	Portainer    ServiceConfig    `koanf:"portainer"`
	Grafana      ServiceConfig    `koanf:"grafana"`
	Prometheus   ServiceConfig    `koanf:"prometheus"`
	SilverBullet ServiceConfig    `koanf:"silverbullet"`
	Vikunja      ServiceConfig    `koanf:"vikunja"`
	Timeout      TimeoutConfig    `koanf:"timeout"`
	TLS          TLSConfig        `koanf:"tls"`
}

type ServerConfig struct {
	Name            string `koanf:"name"`
	Version         string `koanf:"version"`
	ProtocolVersion string `koanf:"protocol_version"`
	HTTPEnabled     bool   `koanf:"http_enabled"`
	HTTPPort        int    `koanf:"http_port"`
	APIEnabled      bool   `koanf:"api_enabled"`
	APIPort         int    `koanf:"api_port"`
	APIToken        string `koanf:"api_token"`
}

type LogConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
}

type ServiceConfig struct {
	URL     string `koanf:"url"`
	Token   string `koanf:"token"`
	Enabled bool   `koanf:"enabled"`
}

type TimeoutConfig struct {
	HTTP      time.Duration `koanf:"http"`
	Operation time.Duration `koanf:"operation"`
}

type TLSConfig struct {
	SkipVerify bool `koanf:"skip_verify"`
}

// Load loads configuration from files and environment variables
func Load(env string) (*Config, error) {
	k := koanf.New(".")

	// Load base config
	if err := k.Load(file.Provider("config/base.yaml"), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading base config: %w", err)
	}

	// Load environment-specific config
	configFile := fmt.Sprintf("config/%s.yaml", env)
	if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
		// Environment config is optional
		fmt.Printf("Warning: could not load %s: %v\n", configFile, err)
	}

	// Load environment variables with APP_ prefix
	// Use __ as delimiter for nested keys (e.g., APP_PORTAINER__TOKEN)
	if err := k.Load(koanfenv.Provider("APP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "__", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading env vars: %w", err)
	}

	// Unmarshal into config struct
	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}
