package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

func dbConnect() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// envFile, _ := godotenv.Read("../.env")
	// dbname := envFile["db_name"]
	// password := envFile["db_pass"]
	// user := envFile["db_user"]
	// host := envFile["db_host"]
	// port := envFile["db_port"]

	dbname := os.Getenv("db_name")
	password := os.Getenv("db_pass")
	user := os.Getenv("db_user")
	host := os.Getenv("db_host")
	port := os.Getenv("db_port") //this comes in as a string :(

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
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
