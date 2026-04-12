package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/alliemony/neo/backend/internal/config"
	"github.com/alliemony/neo/backend/internal/database"
	"github.com/alliemony/neo/backend/internal/handler"
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

	postRepo := repository.NewPostRepo(db)
	postService := service.NewPostService(postRepo)

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

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.Health)
		r.Get("/posts", postHandler.List)
		r.Get("/posts/{slug}", postHandler.GetBySlug)
		r.Get("/tags", postHandler.ListTags)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("neo backend listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
