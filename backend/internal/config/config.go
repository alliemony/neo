package config

import (
	"os"
	"strings"
)

type Config struct {
	Port             string
	DatabaseURL      string
	AdminUsername    string
	AdminPassword    string
	CORSOrigins      string
	WidgetServiceURL string
	ServeStatic      bool
	StaticDir        string
	// OAuth settings
	AuthMode         string // "basic" or "oauth"
	GitHubClientID   string
	GitHubClientSecret string
	OAuthAllowedUsers []string
	SessionSecret    string
	BaseURL          string
}

func Load() Config {
	allowedUsersStr := getEnv("OAUTH_ALLOWED_USERS", "")
	var allowedUsers []string
	if allowedUsersStr != "" {
		for _, u := range strings.Split(allowedUsersStr, ",") {
			if trimmed := strings.TrimSpace(u); trimmed != "" {
				allowedUsers = append(allowedUsers, trimmed)
			}
		}
	}

	return Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "sqlite://neo.db"),
		AdminUsername:     getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:     getEnv("ADMIN_PASSWORD", "changeme"),
		CORSOrigins:       getEnv("CORS_ORIGINS", "http://localhost:5173"),
		WidgetServiceURL:  getEnv("WIDGET_SERVICE_URL", "http://localhost:8000"),
		ServeStatic:       getEnv("SERVE_STATIC", "") == "true",
		StaticDir:         getEnv("STATIC_DIR", "./static"),
		AuthMode:          getEnv("AUTH_MODE", "basic"),
		GitHubClientID:    getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		OAuthAllowedUsers: allowedUsers,
		SessionSecret:     getEnv("SESSION_SECRET", "change-me-to-a-random-32-byte-key"),
		BaseURL:           getEnv("BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
