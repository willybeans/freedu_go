package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
}

func NewServer(r *chi.Mux) *Server {
	fmt.Println("****Server Started on: ", "8080", "****")
	return &Server{
		httpServer: &http.Server{Addr: ":8080", Handler: r},
		router:     r,
	}
}
