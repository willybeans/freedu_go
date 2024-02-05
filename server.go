package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/willybeans/freedu_go/logger"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
}

func NewServer(r *chi.Mux) *Server {
	// fmt.Println("****Server Started on: ", "8080", "****")
	l := logger.Get()
	port := os.Getenv("PORT")
	// addr := fmt.Sprintf(":%s", port)
	l.Info().
		Str("port", port).
		Msgf("Starting App Server on Port '%s'", port)
	return &Server{
		httpServer: &http.Server{Addr: ":" + port, Handler: r},
		router:     r,
	}
}
