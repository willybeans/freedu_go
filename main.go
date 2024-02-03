package main

import (
	"log"

	"github.com/willybeans/freedu_go/database"

	_ "github.com/lib/pq"
)

func main() {

	// Init Router
	router := NewRouter()
	// Init Database
	database.DbConnect()
	defer database.CloseDB()
	s := NewServer(router)
	// Init Server
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
