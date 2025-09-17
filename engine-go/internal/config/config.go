package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds runtime configuration for the Go engine service.
// Only include what we need for initial health endpoints; extend as needed.
type Config struct {
	// HTTP
	HTTPPort int    `envconfig:"PORT" default:"8080"`
	BaseURL  string `envconfig:"BASE_URL"`

	// Maintenance
	CurrentlyUnderMaintenanceMessage string `envconfig:"CURRENTLY_UNDERGOING_MAINTENANCE_MESSAGE"`

	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// Startup probe specific
	// CacheWarmup toggles whether to attempt any lightweight warmups during startup probe
	CacheWarmup bool          `envconfig:"CACHE_WARMUP_ENABLED" default:"true"`
	WarmupTimeout time.Duration `envconfig:"WARMUP_TIMEOUT" default:"2s"`
}

// Load reads the environment into Config.
func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}