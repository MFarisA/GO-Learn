package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"example.com/go-rest/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type NoteHandler struct {
	DB *gorm.DB
}

func NewNoteHandler(db *gorm.DB) *NoteHandler {
	return &NoteHandler{DB: db}
}

func (handler *NoteHandler) CreateNote(writer http.ResponseWriter, request *http.Request) {
	var NewNote models.Note

	err := json.NewDecoder(request.Body).Decode(&NewNote)
	if err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if NewNote.Title == "" || NewNote.Content == "" {
		http.Error(writer, "Title and content must not be empty", http.StatusBadRequest)
		return
	}

	NewNote.CreatedAt = time.Now()
	NewNote.UpdatedAt = time.Now()

	if result := handler.DB.Create(&NewNote); result.Error != nil {
		http.Error(writer, "Failed to create note"+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(NewNote)
}

func (handler *NoteHandler) GetNotes(writer http.ResponseWriter, request *http.Request) {
	var note []models.Note

	if result := handler.DB.Find(&note); result.Error != nil {
		http.Error(writer, "Failed to fetch notes"+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(note)
}

func (handler *NoteHandler) ShowNotesByID(writer http.ResponseWriter, request *http.Request) {
	IdParam := chi.URLParam(request, "id")
	id, err := strconv.Atoi(IdParam)

	if err != nil {
		http.Error(writer, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var Note models.Note

	if result := handler.DB.First(&Note, id); result.Error != nil {
		http.Error(writer, "Failed to fetch note"+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(Note)
}

func (handler *NoteHandler) UpdateNotes(writer http.ResponseWriter, request *http.Request) {
	IdParam := chi.URLParam(request, "id")
	id, err := strconv.Atoi(IdParam)

	if err != nil {
		http.Error(writer, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var updateNote models.Note

	err = json.NewDecoder(request.Body).Decode(&updateNote)
	if err != nil {
		http.Error(writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateNote.Title == "" || updateNote.Content == "" {
		http.Error(writer, "Title and content must not be empty", http.StatusBadRequest)
		return
	}

	var existingNote models.Note
	if result := handler.DB.First(&existingNote, id); result.Error != nil {
		http.Error(writer, "Failed to find note: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	existingNote.Title = updateNote.Title
	existingNote.Content = updateNote.Content
	existingNote.UpdatedAt = existingNote.UpdatedAt

	if result := handler.DB.Save(&existingNote); result.Error != nil {
		http.Error(writer, "Failed to update note: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(existingNote)
}

func (handler *NoteHandler) DeleteNotes(writer http.ResponseWriter, request *http.Request) {
	IdParam := chi.URLParam(request, "id")
	id, err := strconv.Atoi(IdParam)

	if err != nil {
		http.Error(writer, "Invalid note ID", http.StatusBadRequest)
		return
	}

	result := handler.DB.Delete(&models.Note{}, id)
	if result.Error != nil {
		http.Error(writer, "Failed to delete note: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.NotFound(writer, request)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	responseMessage := map[string]string{
		"message": "Note deleted successfully",
	}
	json.NewEncoder(writer).Encode(responseMessage)
}
