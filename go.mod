module github.com/casapps/casspeed

go 1.24.0

require (
	// Network/HTTP
	github.com/go-chi/chi/v5 v5.2.0 // Router
	github.com/go-chi/cors v1.2.1 // CORS (chi-compatible)
	github.com/google/uuid v1.6.0 // UUID generation
	github.com/gorilla/websocket v1.5.3 // WebSocket

	// Utilities
	github.com/robfig/cron/v3 v3.0.1 // Scheduler

	// Core
	gopkg.in/yaml.v3 v3.0.1 // YAML config
	// Database drivers
	modernc.org/sqlite v1.34.5 // SQLite (pure Go)
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	modernc.org/libc v1.55.3 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
)
