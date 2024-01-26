package database

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

func DbConnect() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT") //this comes in as a string :(

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println("Could not connect to DB!")
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	database = db
	log.Println("Connected to DB.")
	// return db

}

func DB() *sql.DB {
	return database
}

func CloseDB() error {
	return database.Close()
}
