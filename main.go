package main

import (
	"github.com/willybeans/freedu_go/database"
	"github.com/willybeans/freedu_go/logger"

	_ "github.com/lib/pq"
)

func main() {
	l := logger.Get()

	// Init Router
	router := NewRouter()
	// Init Database
	database.DbConnect()
	defer database.CloseDB()
	s := NewServer(router)
	// Init Server
	err := s.httpServer.ListenAndServe()
	if err != nil {
		l.Error().Err(err).Msg("server error")
		l.Fatal().Msg("App Server Closed")
	}
}
