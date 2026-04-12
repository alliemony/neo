package config

import "os"

type Config struct {
	Port            string
	DatabaseURL     string
	AdminUsername   string
	AdminPassword   string
	CORSOrigins     string
	WidgetServiceURL string
	ServeStatic     bool
	StaticDir       string
}

func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "sqlite://neo.db"),
		AdminUsername:   getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:   getEnv("ADMIN_PASSWORD", "changeme"),
		CORSOrigins:     getEnv("CORS_ORIGINS", "http://localhost:5173"),
		WidgetServiceURL: getEnv("WIDGET_SERVICE_URL", "http://localhost:8000"),
		ServeStatic:     getEnv("SERVE_STATIC", "") == "true",
		StaticDir:       getEnv("STATIC_DIR", "./static"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
