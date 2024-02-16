package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/willybeans/freedu_go/database"
)

type User_ID struct {
	ID int `json:"id"`
}

type UpdateUser struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
	Profile  string `json:"profile"`
}

type NewUser struct {
	UserName string `json:"username"`
}

type User struct {
	ID          string `json:"id"`
	UserName    string `json:"username"`
	Profile     string `json:"profile"`
	TimeCreated string `json:"time_created"`
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")

	var user User
	err := database.DB().QueryRow("SELECT * FROM users WHERE id = $1", userId).Scan(&user.ID, &user.UserName, &user.Profile, &user.TimeCreated)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB().Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	userList := make([]User, 0)
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.UserName, &user.Profile, &user.TimeCreated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userList = append(userList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

func NewUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser NewUser
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	err := database.DB().QueryRow("INSERT INTO users (username) VALUES ($1) RETURNING *", newUser.UserName).Scan(&user.ID, &user.UserName, &user.TimeCreated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var updateUser UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user User
	query := database.DB().QueryRow("UPDATE users SET username = $1, profile = $2 WHERE id = $3 RETURNING *", updateUser.UserName, updateUser.Profile, updateUser.ID)
	err := query.Scan(&user.ID, &user.UserName, &user.Profile, &user.TimeCreated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var deleteUser ID
	if err := json.NewDecoder(r.Body).Decode(&deleteUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//future: add some user validation
	result, err := database.DB().Exec("DELETE FROM users WHERE id = $1", deleteUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		// panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}
