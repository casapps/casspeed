package admin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/casapps/casspeed/src/server/model"
	"github.com/casapps/casspeed/src/server/store"
	"golang.org/x/crypto/argon2"
)

type Handler struct {
	store store.Store
}

func NewHandler(st store.Store) *Handler {
	return &Handler{store: st}
}

func HashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hex.EncodeToString(salt) + "$" + hex.EncodeToString(hash)
}

func VerifyPassword(password, stored string) bool {
	parts := []byte(stored)
	if len(parts) < 33 {
		return false
	}
	saltHex := string(parts[:32])
	hashHex := string(parts[33:])
	
	salt, _ := hex.DecodeString(saltHex)
	storedHash, _ := hex.DecodeString(hashHex)
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	
	if len(hash) != len(storedHash) {
		return false
	}
	for i := range hash {
		if hash[i] != storedHash[i] {
			return false
		}
	}
	return true
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "web/templates/admin/login.html")
		return
	}

	ctx := r.Context()
	username := r.FormValue("username")
	password := r.FormValue("password")

	admin, err := h.store.GetAdminByUsername(ctx, username)
	if err != nil || admin == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !admin.LockedUntil.IsZero() && time.Now().Before(admin.LockedUntil) {
		http.Error(w, "Account locked. Try again later.", http.StatusForbidden)
		return
	}

	if !VerifyPassword(password, admin.Password) {
		attempts := admin.FailedAttempts + 1
		h.store.UpdateAdminFailedAttempts(ctx, admin.ID, attempts)
		
		if attempts >= 5 {
			h.store.LockAdmin(ctx, admin.ID, time.Now().Add(15*time.Minute))
			http.Error(w, "Too many failed attempts. Account locked for 15 minutes.", http.StatusForbidden)
			return
		}
		
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	h.store.UpdateAdminLastLogin(ctx, admin.ID)

	sessionID := make([]byte, 32)
	rand.Read(sessionID)
	sessionIDHex := hex.EncodeToString(sessionID)

	session := &model.AdminSession{
		ID:        sessionIDHex,
		AdminID:   admin.ID,
		IPAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := h.store.CreateAdminSession(ctx, session); err != nil {
		http.Error(w, "Session creation failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    sessionIDHex,
		Path:     "/admin",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	cookie, err := r.Cookie("admin_session")
	if err == nil {
		h.store.DeleteAdminSession(ctx, cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/admin",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/admin/dashboard.html")
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings := map[string]interface{}{
		"server": map[string]interface{}{
			"port": 80,
			"mode": "production",
		},
		"test": map[string]interface{}{
			"max_threads":      16,
			"default_duration": 10,
			"max_concurrent":   3,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"status":  "success",
		"message": "Settings updated",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		
		cookie, err := r.Cookie("admin_session")
		if err != nil {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		session, err := h.store.GetAdminSession(ctx, cookie.Value)
		if err != nil || session == nil {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}

		h.store.UpdateAdminSessionActivity(ctx, session.ID)

		ctx = context.WithValue(ctx, "admin_id", session.AdminID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// ServerSettings shows settings page
func (h *Handler) ServerSettings(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/admin/settings.html")
}

// ServerInfo shows server info page
func (h *Handler) ServerInfo(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/admin/info.html")
}

// ServerLogs shows logs page
func (h *Handler) ServerLogs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/admin/logs.html")
}
