package swagger

import (
	"net/http"
	"strings"
)

// Theme represents a UI theme
type Theme string

const (
	ThemeLight Theme = "light"
	ThemeDark  Theme = "dark"
	ThemeAuto  Theme = "auto"
)

// DetectTheme detects theme from request
func DetectTheme(r *http.Request) Theme {
	// Check query parameter
	if theme := r.URL.Query().Get("theme"); theme != "" {
		switch theme {
		case "light":
			return ThemeLight
		case "dark":
			return ThemeDark
		case "auto":
			return ThemeAuto
		}
	}
	
	// Check cookie
	if cookie, err := r.Cookie("theme"); err == nil {
		switch cookie.Value {
		case "light":
			return ThemeLight
		case "dark":
			return ThemeDark
		case "auto":
			return ThemeAuto
		}
	}
	
	// Check prefers-color-scheme from headers (user-agent hints)
	// Default to dark for consistency
	return ThemeDark
}

// GetThemeCSS returns CSS for the specified theme
func GetThemeCSS(theme Theme) string {
	switch theme {
	case ThemeLight:
		return `
			body {
				background-color: #ffffff;
				color: #000000;
			}
			.swagger-ui .topbar {
				background-color: #f0f0f0;
			}
		`
	case ThemeDark:
		return `
			body {
				background-color: #1a1a1a;
				color: #ffffff;
			}
			.swagger-ui .topbar {
				background-color: #2a2a2a;
			}
			.swagger-ui {
				filter: invert(88%) hue-rotate(180deg);
			}
			.swagger-ui .renderedMarkdown code,
			.swagger-ui .response .microlight {
				filter: invert(100%) hue-rotate(180deg);
			}
		`
	case ThemeAuto:
		return `
			@media (prefers-color-scheme: dark) {
				body {
					background-color: #1a1a1a;
					color: #ffffff;
				}
				.swagger-ui .topbar {
					background-color: #2a2a2a;
				}
				.swagger-ui {
					filter: invert(88%) hue-rotate(180deg);
				}
				.swagger-ui .renderedMarkdown code,
				.swagger-ui .response .microlight {
					filter: invert(100%) hue-rotate(180deg);
				}
			}
			@media (prefers-color-scheme: light) {
				body {
					background-color: #ffffff;
					color: #000000;
				}
				.swagger-ui .topbar {
					background-color: #f0f0f0;
				}
			}
		`
	default:
		return ""
	}
}

// SetThemeCookie sets theme cookie
func SetThemeCookie(w http.ResponseWriter, theme Theme) {
	cookie := &http.Cookie{
		Name:     "theme",
		Value:    string(theme),
		Path:     "/",
		MaxAge:   31536000,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// ThemeFromString converts string to Theme
func ThemeFromString(s string) Theme {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "light":
		return ThemeLight
	case "dark":
		return ThemeDark
	case "auto":
		return ThemeAuto
	default:
		return ThemeDark
	}
}
