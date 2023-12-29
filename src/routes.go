package main

import (
	"api/internal/handlers"

	"github.com/go-chi/chi"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/healthcheck", handlers.Healthcheck)
	router.Post("/login", handlers.LoginHandler)
	router.Post("/image", handlers.ImageHandler)
	router.Get("/scrape", handlers.ScrapeHandler)

	// s.router.HandleFunc("/api/", s.handleAPI())
	// s.router.HandleFunc("/about", s.handleAbout())
	// s.router.HandleFunc("/", s.handleIndex())

	return router
}
