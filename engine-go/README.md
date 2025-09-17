# engine-go (experimental)

A minimal Go port scaffold for Grafana OnCall engine HTTP endpoints.

Status: Experimental. Provides health, readiness, startup probe, and maintenance status routes.

## Quickstart

- Build:

```bash
make build
```

- Run:

```bash
PORT=8080 CURRENTLY_UNDERGOING_MAINTENANCE_MESSAGE="maintenance" make run
```

- Endpoints:

- `/`, `/health/`, `/api/internal/v1/health/` -> `Ok`
- `/ready/` -> `Ok`
- `/startupprobe/` -> `Ok` (performs optional warmups)
- `/api/internal/v1/maintenance-mode-status` -> `{ "currently_undergoing_maintenance_message": "..." }`

## Configuration

Environment variables:

- `PORT` (default `8080`)
- `CURRENTLY_UNDERGOING_MAINTENANCE_MESSAGE`
- `BASE_URL`
- `LOG_LEVEL` (default `info`)
- `CACHE_WARMUP_ENABLED` (default `true`)
- `WARMUP_TIMEOUT` (default `2s`)

## Build flags

The binary embeds a version string from git using `-ldflags`.