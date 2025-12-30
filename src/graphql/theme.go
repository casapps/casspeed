package graphql

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
	
	// Default to dark
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
			.graphiql-container {
				--color-primary: 40, 130, 250;
				--color-secondary: 110, 110, 110;
			}
		`
	case ThemeDark:
		return `
			body {
				background-color: #1a1a1a;
				color: #ffffff;
			}
			.graphiql-container {
				--color-primary: 100, 180, 255;
				--color-secondary: 180, 180, 180;
				--color-base: 26, 26, 26;
			}
		`
	case ThemeAuto:
		return `
			@media (prefers-color-scheme: dark) {
				body {
					background-color: #1a1a1a;
					color: #ffffff;
				}
				.graphiql-container {
					--color-primary: 100, 180, 255;
					--color-secondary: 180, 180, 180;
					--color-base: 26, 26, 26;
				}
			}
			@media (prefers-color-scheme: light) {
				body {
					background-color: #ffffff;
					color: #000000;
				}
				.graphiql-container {
					--color-primary: 40, 130, 250;
					--color-secondary: 110, 110, 110;
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
