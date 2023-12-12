package main

import (
	"api/internal/handlers"

	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/healthcheck", handlers.Healthcheck)
	// s.router.HandleFunc("/api/", s.handleAPI())
	// s.router.HandleFunc("/about", s.handleAbout())
	// s.router.HandleFunc("/", s.handleIndex())

	// // API routes
	// router.Get("/healthcheck", healthcheck)
	// router.Get("/scrape", scrape)
	// router.Post("/register", registerHandler)
	// router.Post("/login", loginHandler)
	return router
}
