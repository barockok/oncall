package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

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

	// Health endpoints
	r.Get("/", okHandler)
	r.Get("/health/", okHandler)
	r.Get("/api/internal/v1/health/", okHandler)
	r.Get("/ready/", okHandler)

	// Startup probe
	r.Get("/startupprobe/", func(w http.ResponseWriter, r *http.Request) {
		if cfg.CacheWarmup {
			deadline := time.Now().Add(cfg.WarmupTimeout)
			// Placeholder: warmup tasks can be added here, e.g., cache priming, DB ping.
			log.Ctx(r.Context()).Debug().Msg("startup warmup completed")
			_ = deadline
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Ok"))
	})

	// Maintenance mode status
	r.Get("/api/internal/v1/maintenance-mode-status", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"currently_undergoing_maintenance_message": cfg.CurrentlyUnderMaintenanceMessage,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	return r
}

func okHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}