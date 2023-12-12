package main

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	ID        int    // `json:"id"` //was int
	Username  string // `json:"username"`
	Password  string // `json:"-"`
	Time      int
	Firstname string
	Lastname  string
	Email     string
	Age       int
}

var database *sql.DB

// change to .env
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "lang_api"
)

// Connect to the "world" database
func dbConnect() {
	// replace "root" and "password" with your database login credentials
	// db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/world")
	// if err != nil {
	// 	log.Println("Could not connect!")
	// }
	// database = db
	// log.Println("Connected.")

	//////////////////

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Connect to the PostgreSQL database

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Could not connect to DB!")
		log.Fatal(err)
	}
	defer db.Close() // defer pushes function call onto list, which is called after the surrounding function is complete.
	// this is commonly used to simply functions that perform various cleanup tasks, ie, closing the db here

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	database = db
	log.Println("Connected to DB.")
}

func DB() *sql.DB {
	return database
}
