module github.com/casapps/casspeed

go 1.23

require (
	// Database drivers
	modernc.org/sqlite v1.34.5 // SQLite (pure Go)

	// Cache/Cluster
	github.com/redis/go-redis/v9 v9.7.0 // Valkey/Redis

	// Core
	gopkg.in/yaml.v3 v3.0.1           // YAML config
	github.com/google/uuid v1.6.0     // UUID generation
	golang.org/x/crypto v0.31.0       // Argon2, Bcrypt

	// Authentication (required for admin panel)
	github.com/pquerna/otp v1.4.0                   // TOTP 2FA
	github.com/go-webauthn/webauthn v0.11.2         // Passkeys/WebAuthn
	github.com/golang-jwt/jwt/v5 v5.2.1             // JWT tokens
	github.com/coreos/go-oidc/v3 v3.11.0            // OIDC client
	golang.org/x/oauth2 v0.24.0                     // OAuth2 flows
	github.com/go-ldap/ldap/v3 v3.4.10              // LDAP/AD
	github.com/gorilla/sessions v1.4.0              // Cookie sessions

	// Network/HTTP
	github.com/go-chi/chi/v5 v5.2.0       // Router
	github.com/cretz/bine v0.2.0          // Tor controller
	github.com/gorilla/websocket v1.5.3   // WebSocket
	github.com/rs/cors v1.11.1            // CORS middleware

	// Utilities
	github.com/robfig/cron/v3 v3.0.1                // Scheduler
	golang.org/x/time v0.8.0                        // Rate limiting
	github.com/go-playground/validator/v10 v10.23.0 // Validation
)
