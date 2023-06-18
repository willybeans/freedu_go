package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int // `json:"id"` //was int
	Username string // `json:"username"`
	Password string // `json:"-"`
	Time 		 int
	Firstname string
	Lastname string
	Email string
	Age int
}

var db *sql.DB

const (
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "password"
  dbname   = "lang_api"
)

func main() {
	// Initialize the Chi router
	router := chi.NewRouter()
	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// API routes
	router.Get("/healthcheck", healthcheck)
	router.Post("/register", registerHandler)
	router.Post("/login", loginHandler)
	
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	// Connect to the PostgreSQL database
	var err error

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // defer pushes function call onto list, which is called after the surrounding function is complete. 
	// this is commonly used to simply functions that perform various cleanup tasks, ie, closing the db here 

	// Create a new context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	Ping(ctx)

	// Create an HTTP server with the Chi router
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start the server in a separate goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Println("Server started")

	// Wait for interrupt or termination signal to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Call cancel to stop ongoing requests and cleanup resources
	cancel()

	// Create a deadline to wait for outstanding requests to complete before shutting down
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	// Shut down the server gracefully with the specified timeout
	err = server.Shutdown(ctxShutdown)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server stopped")
}

// Ping the database to verify DSN provided by the user is valid and the
// server accessible. If the ping fails exit the program with an error.
func Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Healthcheck!"})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode the request body into the User struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash the user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// get ID and time
	var id int;

	rows, err := db.Query("SELECT nextval('id_seq');")
	if err != nil {
		s := err.Error()
		fmt.Printf("Failed to retreive sequence ID for user id\n s: %v", s)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	
	now := time.Now().Unix()

	// Save the user to the database
	_, err = db.Exec("INSERT INTO users (username, password, time_created, id, age, first_name, last_name, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", user.Username, hashedPassword, now, id, user.Age, user.Firstname, user.Lastname, user.Email)
	if err != nil {
		s := err.Error()
    fmt.Printf("type: %T; value: %q\n", s, s)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode the request body into the User struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve the user from the database
	row := db.QueryRow("SELECT * FROM users WHERE username = $1", user.Username)
	err = row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
