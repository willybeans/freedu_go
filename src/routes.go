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

	router.Get("/getContent", handlers.GetContentHandler)
	router.Get("/getAllContent", handlers.GetAllContentHandler)
	router.Post("/makeNewContent", handlers.NewContentHandler)
	router.Put("/updateContent", handlers.UpdateContentHandler)
	router.Delete("/deleteContent", handlers.DeleteContentHandler)

	// Middleware
	// router.Use(middleware.Logger)
	// router.Use(middleware.Recoverer)
	return router
}
