package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/barockok/oncall/engine-go/internal/config"
)

func notImplementedHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		_, _ = w.Write([]byte(name + " endpoint is not implemented in engine-go yet"))
	}
}

func IntegrationsRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("integrations")))
	return r
}

func GrafanaIncidentRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("api-gi")))
	return r
}

func InternalAPIRouterWithConfig(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	// /features
	r.Get("/features", func(w http.ResponseWriter, r *http.Request) {
		features := computeEnabledFeatures(cfg, r)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(features)
	})

	// default catch-all as 501 for other internal routes
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("api-internal")))
	return r
}

func InternalAPIRouter() http.Handler { // legacy helper for stubs
	return InternalAPIRouterWithConfig(config.Config{})
}

func PluginRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("grafana-plugin")))
	return r
}

func TwilioAppRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("twilioapp")))
	return r
}

func PublicAPIRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("api-public")))
	return r
}

func MobileAppRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("mobile_app")))
	return r
}

func MetricsRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", promhttp.Handler())
	return r
}

func TelegramRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("telegram")))
	return r
}

func SlackRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("slack")))
	return r
}

func MattermostRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("mattermost")))
	return r
}

func OSSInstallationRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("oss_installation")))
	return r
}

func ZvonokRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("zvonok")))
	return r
}

func ExotelRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("exotel")))
	return r
}

func ChatOpsProxyRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("chatops_proxy")))
	return r
}

func DebugToolbarRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("debug_toolbar")))
	return r
}

func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("django_admin")))
	return r
}

func SilkRouter() http.Handler {
	r := chi.NewRouter()
	r.Method(http.MethodGet, "/*", http.HandlerFunc(notImplementedHandler("silk")))
	return r
}

// computeEnabledFeatures mirrors Django FeaturesAPIView logic using env-driven config.
func computeEnabledFeatures(cfg config.Config, r *http.Request) []string {
	var enabled []string

	if cfg.FeatureSlackIntegration {
		enabled = append(enabled, "slack")
	}
	if cfg.UnifiedSlackAppEnabled {
		enabled = append(enabled, "unified_slack")
	}
	if cfg.FeatureTelegramIntegration {
		enabled = append(enabled, "telegram")
	}

	if cfg.IsOpenSource {
		enabled = append(enabled, "grafana_cloud_connection")
		if getBoolHeader(r, "X-Feature-Live-Settings", true) { // proxy for live_settings.GRAFANA_CLOUD_NOTIFICATIONS_ENABLED env
			if getBoolHeader(r, "X-Feature-Live-Settings-Enabled", true) && getBoolHeader(r, "X-Feature-Live-Settings-Flag", true) {
				// This is a placeholder; real implementation would consult live settings
			}
		}
		// Allow enabling live settings via env or header override
		if getBoolHeader(r, "X-Feature-Live-Settings", true) {
			enabled = append(enabled, "live_settings")
		}
		if getBoolHeader(r, "X-Feature-Grafana-Cloud-Notifications", false) {
			enabled = append(enabled, "grafana_cloud_notifications")
		}
	} else {
		enabled = append(enabled, "msteams")
	}

	if getBoolHeader(r, "X-Feature-Grafana-Alerting-V2", true) {
		enabled = append(enabled, "grafana_alerting_v2")
	}

	// Labels feature: without full org and DB, gate via header to simulate org capabilities
	if getBoolHeader(r, "X-Org-Labels-Enabled", false) {
		enabled = append(enabled, "labels")
	}

	if getBoolHeader(r, "X-Feature-Google-OAuth2", false) {
		enabled = append(enabled, "google_oauth2")
	}
	if getBoolHeader(r, "X-Feature-Service-Dependencies", false) {
		enabled = append(enabled, "service_dependencies")
	}
	if getBoolHeader(r, "X-Feature-Personal-Webhook", true) {
		enabled = append(enabled, "personal_webhook")
	}
	if cfg.FeatureMattermostIntegration {
		enabled = append(enabled, "mattermost")
	}

	return enabled
}

func getBoolHeader(r *http.Request, key string, def bool) bool {
	v := strings.ToLower(strings.TrimSpace(r.Header.Get(key)))
	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	}
	return def
}