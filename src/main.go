package main

import (
	"api/internal/database"
	"log"

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

// func registerHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User

// 	// Decode the request body into the User struct
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	// Hash the user password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
// 		return
// 	}

// 	// get ID and time
// 	var id int

// 	rows, err := db.Query("SELECT nextval('id_seq');")
// 	if err != nil {
// 		s := err.Error()
// 		fmt.Printf("Failed to retreive sequence ID for user id\n s: %v", s)
// 		http.Error(w, "Failed to register user", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		err := rows.Scan(&id)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		log.Println(id)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	now := time.Now().Unix()

// 	// Save the user to the database
// 	_, err = db.Exec("INSERT INTO users (username, password, time_created, id, age, first_name, last_name, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", user.Username, hashedPassword, now, id, user.Age, user.Firstname, user.Lastname, user.Email)
// 	if err != nil {
// 		s := err.Error()
// 		fmt.Printf("type: %T; value: %q\n", s, s)
// 		http.Error(w, "Failed to register user", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User

// 	// Decode the request body into the User struct
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve the user from the database
// 	row := db.QueryRow("SELECT * FROM users WHERE username = $1", user.Username)
// 	err = row.Scan(&user.ID, &user.Username, &user.Password)
// 	if err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	// Compare the provided password with the hashed password
// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Password))
// 	if err != nil {
// 		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
// }
