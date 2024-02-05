package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/willybeans/freedu_go/logger"
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

	l := logger.Get()

	err := godotenv.Load(".env")
	if err != nil {
		l.Fatal().
			Msg(" Error Loading .env file ")
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
		l.Fatal().
			Msg("Error connecting to DB!")
	}

	err = db.Ping()
	if err != nil {
		l.Fatal().
			Msg("Error on DB ping!")
	}

	database = db
	l.Info().
		Msgf("Connected to DB on port '%s'", port)
}

func DB() *sql.DB {
	return database
}

func CloseDB() error {
	l := logger.Get()

	l.Fatal().
		Msg("Database closed!")
	return database.Close()
}
