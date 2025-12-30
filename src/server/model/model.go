package model

import "time"

type User struct {
	ID                string    `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordHash      string    `json:"-"`
	ShareShowUsername bool      `json:"share_show_username"`
	CreatedAt         time.Time `json:"created_at"`
}

type Device struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
}

type SpeedTest struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id,omitempty"`
	DeviceID     string    `json:"device_id,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	DownloadMbps float64   `json:"download_mbps"`
	UploadMbps   float64   `json:"upload_mbps"`
	PingMs       float64   `json:"ping_ms"`
	JitterMs     float64   `json:"jitter_ms"`
	PacketLoss   float64   `json:"packet_loss"`
	ClientIPHash string    `json:"-"`
	UserAgent    string    `json:"user_agent"`
	ServerID     string    `json:"server_id"`
	ShareCode    string    `json:"share_code,omitempty"`
	ShareViews   int       `json:"share_views"`
	CreatedAt    time.Time `json:"created_at"`
}

type APIToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"-"`
	Name      string    `json:"name"`
	LastUsed  time.Time `json:"last_used,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Data      string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Admin struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Password       string    `json:"-"`
	Email          string    `json:"email,omitempty"`
	Role           string    `json:"role"`
	Enabled        bool      `json:"enabled"`
	APITokenHash   string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastLogin      time.Time `json:"last_login,omitempty"`
	FailedAttempts int       `json:"-"`
	LockedUntil    time.Time `json:"-"`
}

type AdminSession struct {
	ID         string    `json:"id"`
	AdminID    int       `json:"admin_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	LastActive time.Time `json:"last_active"`
}
