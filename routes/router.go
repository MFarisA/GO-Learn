package routes

import (
	"example.com/go-rest/handlers" // Menggunakan handler
	"net/http"
	"strings"
)

func SetupRouter(noteHandler *handlers.NoteHandler) *http.ServeMux{
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