package main

import (
	"time"

	"github.com/willybeans/freedu_go/handlers"
	"github.com/willybeans/freedu_go/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() *chi.Mux {

	router := chi.NewRouter()
	//"github.com/rs/cors"
	// router.Use(cors.Default().Handler)
	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	// router.Use(middleware.Logger(Logger))
	// This middleware will catch and treat panics
	// router.Use(middleware.Recoverer(Logger))

	// Middleware
	router.Use(logger.RequestLogger)
	// router.Use(middleware.Recoverer(logger.RequestLogger))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/healthcheck", handlers.Healthcheck)
	router.Post("/image", handlers.ImageHandler)
	router.Get("/scrape", handlers.ScrapeHandler)

	router.Get("/getContent", handlers.GetContentHandler)
	router.Get("/getAllContent", handlers.GetAllContentHandler)
	router.Post("/newContent", handlers.NewContentHandler)
	router.Put("/updateContent", handlers.UpdateContentHandler)
	router.Delete("/deleteContent", handlers.DeleteContentHandler)

	router.Get("/getUser", handlers.GetUserHandler)
	router.Get("/getAllUsers", handlers.GetAllUsersHandler)
	router.Post("/newUser", handlers.NewUserHandler)
	router.Put("/updateUser", handlers.UpdateUserHandler)
	router.Delete("/deleteUser", handlers.DeleteUserHandler)

	return router
}
