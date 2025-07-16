package handlers

import (
	"database/sql"
	"encoding/json"
	"example.com/go-rest/models" 
	"net/http"
	"strconv"
	"strings"
)

type NoteHandler struct {
	DB *sql.DB
}

type NoteRequest struct{

}

func (handler *NoteHandler) CreateNote(writer http.ResponseWriter, request *http.Request){
	title := request.FormValue("title")
	content := request.FormValue("content")

	if title == "" || content == "" {
		http.Error(writer, "Title and content must not be empty", http.StatusBadRequest)
		return
	}

	result, err := handler.DB.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", title, content)
	if err != nil {
		http.Error(writer, "Failed to create note", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	var newNote models.Note
	err = handler.DB.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?", id).Scan(
		&newNote.ID, &newNote.Title, &newNote.Content, &newNote.CreatedAt, &newNote.UpdatedAt,
	)
	if err != nil {
		http.Error(writer, "Failed to retrieve newly created note", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(newNote)	
}

func (handler *NoteHandler) getNotes(writer http.ResponseWriter, request *http.Request){
	rows, err :=  handler.DB.Query("SELECT id, title, content, created_at, updated_at FROM notes")

	if err != nil {
		http.Error(writer, "Failed to show Notes", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	note := []models.Note{}
}



