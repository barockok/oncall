package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds runtime configuration for the Go engine service.
// Only include what we need for initial endpoints and routing; extend as needed.
type Config struct {
	// HTTP
	HTTPPort int    `envconfig:"PORT" default:"8080"`
	BaseURL  string `envconfig:"BASE_URL"`

	// Maintenance
	CurrentlyUnderMaintenanceMessage string `envconfig:"CURRENTLY_UNDERGOING_MAINTENANCE_MESSAGE"`

	// Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	// Startup probe specific
	CacheWarmup   bool          `envconfig:"CACHE_WARMUP_ENABLED" default:"true"`
	WarmupTimeout time.Duration `envconfig:"WARMUP_TIMEOUT" default:"2s"`

	// Feature flags (mirroring Django settings)
	DetachedIntegrationsServer     bool `envconfig:"DETACHED_INTEGRATIONS_SERVER" default:"false"`
	FeaturePrometheusExporter      bool `envconfig:"FEATURE_PROMETHEUS_EXPORTER_ENABLED" default:"false"`
	FeatureTelegramIntegration     bool `envconfig:"FEATURE_TELEGRAM_INTEGRATION_ENABLED" default:"false"`
	FeatureSlackIntegration        bool `envconfig:"FEATURE_SLACK_INTEGRATION_ENABLED" default:"false"`
	FeatureMattermostIntegration   bool `envconfig:"FEATURE_MATTERMOST_INTEGRATION_ENABLED" default:"false"`
	UnifiedSlackAppEnabled         bool `envconfig:"UNIFIED_SLACK_APP_ENABLED" default:"false"`
	IsOpenSource                   bool `envconfig:"IS_OPEN_SOURCE" default:"true"`
	Debug                          bool `envconfig:"DEBUG" default:"false"`
	SilkProfilerEnabled            bool `envconfig:"SILK_PROFILER_ENABLED" default:"false"`
	DrfSpectacularEnabled          bool `envconfig:"DRF_SPECTACULAR_ENABLED" default:"false"`

	// Paths for admin/silk when enabled
	OncallDjangoAdminPath string `envconfig:"ONCALL_DJANGO_ADMIN_PATH" default:"django-admin/"`
	SilkPath              string `envconfig:"SILK_PATH" default:"silk/"`
}

// Load reads the environment into Config.
func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}