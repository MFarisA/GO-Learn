package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Createdat time.Time `json:"created_at"`
	Updatedat time.Time `json:"updated_at"`
}

var db *sql.DB

func main() {
	var err error

	dsn := "root:@tcp(127.0.0.1:3306)/go_notes_api?parseTime=true"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database : %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	http.HandleFunc("/notes/", notesRouter)

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func notesRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) == 1 && parts[0] == "notes" {
		switch r.Method {
		case http.MethodGet:
			getNotes(w, r)
		case http.MethodPost:
			createNotes(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 2 && parts[0] == "notes" {
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			getNoteByID(w, r, id)
		case http.MethodPut:
			updateNote(w, r, id)
		case http.MethodDelete:
			deleteNote(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func createNotes(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content must not be empty", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", title, content)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	var newNote Note
	err = db.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?", id).Scan(
		&newNote.ID, &newNote.Title, &newNote.Content, &newNote.Createdat, &newNote.Updatedat,
	)
	if err != nil {
		http.Error(w, "Failed to retrieve newly created note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM notes")

	if err != nil {
		http.Error(w, "Failed to Fetch Notes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notes := []Note{}

	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Createdat, &note.Updatedat); err != nil {
			http.Error(w, "Failed to scan note", http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func getNoteByID(w http.ResponseWriter, r *http.Request, id int) {
	var note Note

	err := db.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?", id).Scan(&note.ID, &note.Title, &note.Content, &note.Createdat, &note.Updatedat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Failed to fetch note", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func updateNote(w http.ResponseWriter, r *http.Request, id int) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content must not be empty", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", title, content, id)
	if err != nil {
		http.Error(w, "Failed to update note", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	var updatedNote Note
	err = db.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?", id).Scan(
		&updatedNote.ID, &updatedNote.Title, &updatedNote.Content, &updatedNote.Createdat, &updatedNote.Updatedat,
	)
	if err != nil {
		http.Error(w, "Failed to retrieve updated note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedNote)
}

func deleteNote(w http.ResponseWriter, r *http.Request, id int) {
	result, err := db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete note", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseMessage := map[string]string{
		"message": "Note deleted successfully",
	}

	json.NewEncoder(w).Encode(responseMessage)
}
