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
		// log.Fatal(err)
		l.Fatal().
			// Err(http.ListenAndServe(":"+port, mux)).
			Msg("App Server Closed")
	}
}
