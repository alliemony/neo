package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/alliemony/neo/backend/internal/config"
	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/handler"
	"github.com/alliemony/neo/backend/internal/middleware"
	"github.com/alliemony/neo/backend/internal/repository"
	"github.com/alliemony/neo/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	// Enable foreign key support for cascade deletes.
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Printf("enable foreign keys: %v", err)
	}

	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)
	pageRepo := repository.NewPageRepo(db)

	postService := service.NewPostService(postRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)
	pageService := service.NewPageService(pageRepo)

	if err := database.Seed(db); err != nil {
		log.Printf("seed data: %v", err)
	}

	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.CORSOrigins},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	healthHandler := handler.NewHealthHandler()
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	likeHandler := handler.NewLikeHandler(postRepo)
	pageHandler := handler.NewPageHandler(pageService)
	adminPostHandler := handler.NewAdminPostHandler(postService)
	adminPageHandler := handler.NewAdminPageHandler(pageService)

	commentRateLimiter := middleware.RateLimiter(30 * time.Second)

	siteURL := cfg.CORSOrigins
	rssHandler := handler.NewRSSHandler(postService, siteURL)

	// Select auth middleware based on AUTH_MODE
	var adminAuthMiddleware func(http.Handler) http.Handler
	var oauthAuth *middleware.OAuthAuthenticator

	switch strings.ToLower(cfg.AuthMode) {
	case "oauth":
		oauthAuth = middleware.NewOAuthAuthenticator(middleware.OAuthConfig{
			ClientID:      cfg.GitHubClientID,
			ClientSecret:  cfg.GitHubClientSecret,
			RedirectURL:   cfg.BaseURL + "/api/v1/auth/callback",
			AllowedUsers:  cfg.OAuthAllowedUsers,
			SessionSecret: []byte(cfg.SessionSecret),
			Secure:        cfg.BaseURL != "" && strings.HasPrefix(cfg.BaseURL, "https"),
		})
		adminAuthMiddleware = oauthAuth.Middleware()
		log.Printf("auth mode: oauth (GitHub)")
	default:
		adminAuthMiddleware = middleware.BasicAuth(cfg.AdminUsername, cfg.AdminPassword)
		log.Printf("auth mode: basic")
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.Health)

		// RSS feed
		r.Get("/feed.xml", rssHandler.Feed)

		// Auth mode endpoint (tells frontend which login to show)
		r.Get("/auth/mode", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"mode":"%s"}`, cfg.AuthMode)
		})

		// OAuth auth routes (only active in oauth mode)
		if oauthAuth != nil {
			authHandler := handler.NewAuthHandler(oauthAuth)
			r.Get("/auth/login", authHandler.Login)
			r.Get("/auth/callback", authHandler.Callback)
			r.Get("/auth/me", authHandler.Me)
			r.Post("/auth/logout", authHandler.Logout)
		}

		// Public post routes
		r.Get("/posts", postHandler.List)
		r.Get("/posts/{slug}", postHandler.GetBySlug)
		r.Get("/tags", postHandler.ListTags)

		// Comment routes
		r.Get("/posts/{slug}/comments", commentHandler.List)
		r.With(commentRateLimiter).Post("/posts/{slug}/comments", commentHandler.Create)

		// Like route
		r.Post("/posts/{slug}/like", likeHandler.Like)

		// Public page routes
		r.Get("/pages", pageHandler.ListPublished)
		r.Get("/pages/{slug}", pageHandler.GetBySlug)

		// Admin routes (protected by selected auth middleware)
		r.Route("/admin", func(r chi.Router) {
			r.Use(adminAuthMiddleware)

			r.Get("/posts", adminPostHandler.ListAll)
			r.Post("/posts", adminPostHandler.Create)
			r.Put("/posts/{slug}", adminPostHandler.Update)
			r.Delete("/posts/{slug}", adminPostHandler.Delete)

			r.Get("/pages", adminPageHandler.ListAll)
			r.Post("/pages", adminPageHandler.Create)
			r.Put("/pages/{slug}", adminPageHandler.Update)
			r.Delete("/pages/{slug}", adminPageHandler.Delete)
		})
	})

	// Serve frontend static files in production mode
	if cfg.ServeStatic {
		staticDir := cfg.StaticDir
		fileServer := http.FileServer(http.Dir(staticDir))

		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join(staticDir, r.URL.Path)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// SPA fallback: serve index.html for non-file routes
				http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
				return
			}
			fileServer.ServeHTTP(w, r)
		})
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("neo backend listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
