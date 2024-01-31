package handlers

import (
	"api/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type ID struct {
	ID int `json:"id"`
}

type Author_ID struct {
	ID int `json:"author_id"`
}

type NewContent struct {
	Author_ID   int    `json:"author_id"`
	Title       string `json:"title"`
	BodyContent string `json:"body_content"`
}

type UpdateContent struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	BodyContent string `json:"body_content"`
}

type Content struct {
	Author_ID   int    `json:"author_id"`
	Content_ID  int    `json:"id"`
	Title       string `json:"title"`
	BodyContent string `json:"body_content"`
	TimeCreated string `json:"time_created"`
}

func GetContentHandler(w http.ResponseWriter, r *http.Request) {
	contentId := r.URL.Query().Get("id")

	var content Content
	err := database.DB().QueryRow("SELECT * FROM content WHERE id = $1", contentId).Scan(&content.Content_ID, &content.Author_ID, &content.Title, &content.BodyContent, &content.TimeCreated)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func GetAllContentHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")

	rows, err := database.DB().Query("SELECT * FROM content WHERE author_id=$1", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	contentList := make([]Content, 0)
	for rows.Next() {
		var content Content
		err := rows.Scan(&content.Content_ID, &content.Author_ID, &content.Title, &content.BodyContent, &content.TimeCreated)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contentList = append(contentList, content)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contentList)
}

func NewContentHandler(w http.ResponseWriter, r *http.Request) {
	var newContent NewContent
	if err := json.NewDecoder(r.Body).Decode(&newContent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var content Content
	query := database.DB().QueryRow("INSERT INTO content (title, body_content, author_id) VALUES ($1, $2, $3) RETURNING *", newContent.Title, newContent.BodyContent, newContent.Author_ID)
	err := query.Scan(&content.Content_ID, &content.Author_ID, &content.Title, &content.BodyContent, &content.TimeCreated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(content)
}

func UpdateContentHandler(w http.ResponseWriter, r *http.Request) {
	var updateContent UpdateContent
	if err := json.NewDecoder(r.Body).Decode(&updateContent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// future: add user validation
	var content Content
	err := database.DB().QueryRow("UPDATE content SET title = $1, body_content = $2 WHERE id = $3 RETURNING *", updateContent.Title, updateContent.BodyContent, updateContent.ID).Scan(&content.Content_ID, &content.Author_ID, &content.Title, &content.BodyContent, &content.TimeCreated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(content)

}

func DeleteContentHandler(w http.ResponseWriter, r *http.Request) {
	var deleteContent ID
	if err := json.NewDecoder(r.Body).Decode(&deleteContent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//future: add some user validation
	result, err := database.DB().Exec("DELETE FROM content WHERE id = $1", deleteContent.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	count, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println(count)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}
