package main

import (
	"time"

	"github.com/willybeans/freedu_go/handlers"
	"github.com/willybeans/freedu_go/logger"
	"github.com/willybeans/freedu_go/websockets"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() *chi.Mux {

	router := chi.NewRouter()

	// Middleware
	//"github.com/rs/cors"
	// router.Use(cors.Default().Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.RequestLogger)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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
	router.Get("/getAllContentByQuery", handlers.GetAllContentByQueryHandler)
	router.Post("/newContent", handlers.NewContentHandler)
	router.Put("/updateContent", handlers.UpdateContentHandler)
	router.Delete("/deleteContent", handlers.DeleteContentHandler)

	router.Get("/getUser", handlers.GetUserHandler)
	router.Get("/getAllUsers", handlers.GetAllUsersHandler)
	router.Post("/newUser", handlers.NewUserHandler)
	router.Put("/updateUser", handlers.UpdateUserHandler)
	router.Delete("/deleteUser", handlers.DeleteUserHandler)

	router.Get("/getMessagesByChatId", handlers.GetMessagesByChatIDHandler)
	router.Get("/getAllChatsByUserId", handlers.GetAllChatsForUserHandler)
	// combine with other req for users
	router.Get("/getAllXRefForChat", handlers.GetAllXRefForChatHandler)

	router.Post("/newChat", handlers.NewChatHandler)
	router.Post("/newMessageForUserInChat", handlers.NewMessageForUserInChatHandler)

	router.Handle("/ws", websockets.CreateWsConnection())

	return router
}
