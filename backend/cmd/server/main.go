package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/alliemony/neo/backend/internal/config"
	"github.com/alliemony/neo/backend/internal/handler"
)

func main() {
	cfg := config.Load()

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

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.Health)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("neo backend listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
