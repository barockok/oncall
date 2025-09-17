package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/barockok/oncall/engine-go/internal/config"
)

// NewRouter creates the HTTP router and wires all routes.
func NewRouter(cfg config.Config) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Health endpoints (always available)
	r.Get("/", okHandler)
	r.Get("/health/", okHandler)
	r.Get("/api/internal/v1/health/", okHandler)
	r.Get("/ready/", okHandler)
	r.Get("/startupprobe/", startupProbeHandler(cfg))
	r.Get("/api/internal/v1/maintenance-mode-status", maintenanceHandler(cfg))

	// Integrations
	if !cfg.DetachedIntegrationsServer {
		r.Mount("/integrations/v1/", IntegrationsRouter())
	}

	// Core app routes (stubs)
	r.Mount("/api/gi/v1/", GrafanaIncidentRouter())
	r.Mount("/api/internal/v1/", InternalAPIRouterWithConfig(cfg))
	r.Mount("/api/internal/v1/plugin/", PluginRouter())
	r.Mount("/twilioapp/", TwilioAppRouter())
	r.Mount("/api/v1/", PublicAPIRouter())
	r.Mount("/mobile_app/v1/", MobileAppRouter())
	r.Mount("/api/internal/v1/mobile_app/", MobileAppRouter())

	// Feature conditional routes
	if cfg.FeaturePrometheusExporter {
		r.Mount("/metrics/", MetricsRouter())
	}
	if cfg.FeatureTelegramIntegration {
		r.Mount("/telegram/", TelegramRouter())
	}
	if cfg.FeatureSlackIntegration {
		r.Mount("/api/internal/v1/slack/", SlackRouter())
		r.Mount("/api/v3/webhook/slack/", SlackRouter())
		r.Mount("/slack/", SlackRouter())
	}
	if cfg.FeatureMattermostIntegration {
		r.Mount("/api/internal/v1/mattermost/", MattermostRouter())
	}
	if cfg.IsOpenSource {
		r.Mount("/api/internal/v1/", OSSInstallationRouter())
		r.Mount("/zvonok/", ZvonokRouter())
		r.Mount("/exotel/", ExotelRouter())
	}
	if cfg.UnifiedSlackAppEnabled {
		r.Mount("/api/chatops/", ChatOpsProxyRouter())
	}

	// Debug/Profiler placeholders
	if cfg.Debug {
		r.Mount("/__debug__/", DebugToolbarRouter())
	}
	if cfg.SilkProfilerEnabled {
		r.Mount("/"+cfg.OncallDjangoAdminPath, AdminRouter())
		r.Mount("/"+cfg.SilkPath, SilkRouter())
	}

	// Static files are not handled here; leave to a separate asset server/reverse proxy

	return r
}

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}

func startupProbeHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.CacheWarmup {
			// Placeholder warmup work
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Ok"))
	}
}

func maintenanceHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"currently_undergoing_maintenance_message": cfg.CurrentlyUnderMaintenanceMessage,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}