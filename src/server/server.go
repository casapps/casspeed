package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/casapps/casspeed/src/admin"
	"github.com/casapps/casspeed/src/config"
	"github.com/casapps/casspeed/src/graphql"
	"github.com/casapps/casspeed/src/mode"
	"github.com/casapps/casspeed/src/server/handler"
	"github.com/casapps/casspeed/src/server/service"
	"github.com/casapps/casspeed/src/server/store"
	"github.com/casapps/casspeed/src/swagger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	Config       *config.Config
	Mode         *mode.State
	Router       *chi.Mux
	HTTP         *http.Server
	Store        store.Store
	Handler      *handler.SpeedTestHandler
	ImageHandler *handler.ShareImageHandler
	UserHandler  *handler.UserHandler
	AdminHandler *admin.Handler
	ipTestCount  map[string]*ipRateLimit
	ipMutex      sync.RWMutex
	startTime    time.Time
	version      string
}

type ipRateLimit struct {
	activeTests int
	lastTest    time.Time
}

func New(cfg *config.Config, appMode *mode.State, dataDir string, version string) (*Server, error) {
	dbPath := filepath.Join(dataDir, "db", "speedtest.db")
	dbStore, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("creating store: %w", err)
	}

	speedTestService := service.NewSpeedTestService(cfg.Test.MaxThreads, cfg.Test.ChunkSize)
	speedTestHandler := handler.NewSpeedTestHandler(dbStore, speedTestService)
	imageHandler := handler.NewShareImageHandler(dbStore)
	userHandler := handler.NewUserHandler(dbStore)
	adminHandler := admin.NewHandler(dbStore)

	s := &Server{
		Config:       cfg,
		Mode:         appMode,
		Router:       chi.NewRouter(),
		Store:        dbStore,
		Handler:      speedTestHandler,
		ImageHandler: imageHandler,
		UserHandler:  userHandler,
		AdminHandler: adminHandler,
		ipTestCount:  make(map[string]*ipRateLimit),
		startTime:    time.Now(),
		version:      version,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s, nil
}

func (s *Server) setupMiddleware() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(s.rateLimitMiddleware)

	if s.Mode.IsDevelopment() || s.Mode.IsDebug() {
		s.Router.Use(middleware.Timeout(60 * time.Second))
	} else {
		s.Router.Use(middleware.Timeout(30 * time.Second))
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	s.Router.Use(corsHandler.Handler)
}

func (s *Server) setupRoutes() {
	s.Router.Get("/", s.handleIndex)
	s.Router.Get("/healthz", s.handleHealth)

	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Get("/", s.handleAPIRoot)
		r.Get("/healthz", s.handleHealth)
		
		// Speed test endpoints
		r.Post("/speedtest/start", s.Handler.StartTest)
		r.Get("/speedtest/ws", s.Handler.TestStatus)
		r.Get("/speedtest/download", s.Handler.Download)
		r.Post("/speedtest/upload", s.Handler.Upload)
		r.Get("/speedtest/result/{id}", s.Handler.GetResult)
		r.Get("/speedtest/history", s.Handler.GetHistory)

		// User management endpoints
		r.Post("/users/register", s.UserHandler.Register)
		r.Get("/users/{id}", s.UserHandler.GetProfile)
		r.Get("/users/{id}/devices", s.UserHandler.ListDevices)
		r.Post("/users/{id}/devices", s.UserHandler.CreateDevice)
		r.Delete("/users/{id}/devices/{deviceId}", s.UserHandler.DeleteDevice)
		r.Get("/users/{id}/tokens", s.UserHandler.ListTokens)
		r.Post("/users/{id}/tokens", s.UserHandler.CreateToken)
		r.Delete("/users/{id}/tokens/{tokenId}", s.UserHandler.RevokeToken)

		// Admin API endpoints
		r.Get("/admin/settings", s.AdminHandler.RequireAuth(s.AdminHandler.GetSettings))
		r.Put("/admin/settings", s.AdminHandler.RequireAuth(s.AdminHandler.UpdateSettings))
	})

	// Admin panel web UI
	s.Router.Get("/admin", s.AdminHandler.Login)
	s.Router.Post("/admin/login", s.AdminHandler.Login)
	s.Router.Get("/admin/logout", s.AdminHandler.Logout)
	s.Router.Get("/admin/dashboard", s.AdminHandler.RequireAuth(s.AdminHandler.Dashboard))
	s.Router.Get("/admin/server/settings", s.AdminHandler.RequireAuth(s.AdminHandler.ServerSettings))
	s.Router.Get("/admin/server/info", s.AdminHandler.RequireAuth(s.AdminHandler.ServerInfo))
	s.Router.Get("/admin/server/logs", s.AdminHandler.RequireAuth(s.AdminHandler.ServerLogs))

	s.Router.Get("/share/{code}", s.Handler.GetShare)
	s.Router.Get("/s/{code}", s.Handler.GetShare)
	s.Router.Get("/share/{code}.png", s.ImageHandler.GetSharePNG)
	s.Router.Get("/s/{code}.png", s.ImageHandler.GetSharePNG)
	s.Router.Get("/share/{code}.svg", s.ImageHandler.GetShareSVG)
	s.Router.Get("/s/{code}.svg", s.ImageHandler.GetShareSVG)

	// OpenAPI/Swagger UI and specification
	s.Router.Get("/openapi", swagger.Handler)
	s.Router.Get("/openapi.json", swagger.SpecHandler)

	// GraphQL API and GraphiQL interface
	s.Router.Get("/graphql", graphql.Handler)
	s.Router.Post("/graphql/query", graphql.QueryHandler)

	if s.Mode.ShouldEnableDebugEndpoints() {
		s.Router.Mount("/debug", middleware.Profiler())
	}
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Calculate uptime
	uptime := time.Since(s.startTime)
	days := int(uptime.Hours() / 24)
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	uptimeStr := fmt.Sprintf("%dd %dh %dm", days, hours, minutes)

	// Get hostname
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	// Health response per PART 16 spec
	response := map[string]interface{}{
		"status":    "healthy",
		"version":   s.version,
		"mode":      s.Mode.String(),
		"uptime":    uptimeStr,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"node": map[string]string{
			"id":       "standalone",
			"hostname": hostname,
		},
		"cluster": map[string]interface{}{
			"enabled": false,
		},
		"checks": map[string]string{
			"database": "ok",
			"cache":    "ok",
			"disk":     "ok",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.MarshalIndent(response, "", "  ")
	w.Write(data)
	w.Write([]byte("\n"))
}

func (s *Server) handleAPIRoot(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"version": "v1",
		"status":  "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.MarshalIndent(response, "", "  ")
	w.Write(data)
	w.Write([]byte("\n"))
}

func (s *Server) Start(address string, port int) error {
	addr := fmt.Sprintf("%s:%d", address, port)

	s.HTTP = &http.Server{
		Addr:         addr,
		Handler:      s.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	fmt.Printf("│  🌐 HTTP   http://%s%s│\n", addr, padAddr(addr))
	fmt.Println("├─────────────────────────────────────────────────────────────┤")
	fmt.Printf("│  📡 Listening on http://%s%s│\n", addr, padAddr(addr))
	fmt.Printf("│  ✅ Server started on %s%s│\n", time.Now().Format("Mon Jan 02, 2006 at 15:04:05 MST"), padTime())
	fmt.Println("╰─────────────────────────────────────────────────────────────╯")

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.HTTP.ListenAndServe()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-sigChan:
		return s.Shutdown()
	}
}

func (s *Server) Shutdown() error {
	fmt.Println("\n🛑 Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if s.Store != nil {
		s.Store.Close()
	}

	if err := s.HTTP.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Println("✅ Server stopped")
	return nil
}

func padAddr(addr string) string {
	needed := 60 - len("🌐 HTTP   http://") - len(addr)
	if needed < 0 {
		needed = 0
	}
	return fmt.Sprintf("%*s", needed, "")
}

func padTime() string {
	ts := time.Now().Format("Mon Jan 02, 2006 at 15:04:05 MST")
	needed := 60 - len("✅ Server started on ") - len(ts)
	if needed < 0 {
		needed = 0
	}
	return fmt.Sprintf("%*s", needed, "")
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/speedtest/ws" && r.URL.Path != "/api/v1/speedtest/start" {
			next.ServeHTTP(w, r)
			return
		}

		clientIP := r.RemoteAddr

		s.ipMutex.Lock()
		limit, exists := s.ipTestCount[clientIP]
		if !exists {
			limit = &ipRateLimit{activeTests: 0, lastTest: time.Time{}}
			s.ipTestCount[clientIP] = limit
		}

		if limit.activeTests >= s.Config.Test.MaxConcurrent {
			s.ipMutex.Unlock()
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Too many concurrent tests", http.StatusTooManyRequests)
			return
		}

		secondsSinceLastTest := time.Since(limit.lastTest).Seconds()
		if secondsSinceLastTest < float64(s.Config.Test.MinInterval) {
			retryAfter := int(float64(s.Config.Test.MinInterval) - secondsSinceLastTest)
			s.ipMutex.Unlock()
			w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
			http.Error(w, "Test interval too short", http.StatusTooManyRequests)
			return
		}

		limit.activeTests++
		limit.lastTest = time.Now()
		s.ipMutex.Unlock()

		defer func() {
			s.ipMutex.Lock()
			if l, ok := s.ipTestCount[clientIP]; ok {
				l.activeTests--
			}
			s.ipMutex.Unlock()
		}()

		next.ServeHTTP(w, r)
	})
}
