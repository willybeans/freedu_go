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

	// Middleware
	//"github.com/rs/cors"
	// router.Use(cors.Default().Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.RequestLogger)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/healthcheck", handlers.Healthcheck)
	router.Post("/image", handlers.ImageHandler)
	router.Get("/scrape", handlers.ScrapeHandler)

	router.Get("/getContent", handlers.GetContentHandler)
	router.Get("/getAllContent", handlers.GetAllContentHandler)
	router.Get("/getAllUserContent", handlers.GetAllUserContentHandler)
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
