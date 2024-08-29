package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/blankspace9/notes-app/internal/domain/auth"
	"github.com/blankspace9/notes-app/internal/domain/models"
	"github.com/blankspace9/notes-app/internal/lib/logger/sl"
)

func (h *Handler) addNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.CtxUserID).(int64)
	if !ok {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		h.log.Warn("invalid user id")
		return
	}

	var note models.NoteRequest
	d := json.NewDecoder(r.Body)

	err := d.Decode(&note)
	if err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusBadRequest)
		h.log.Warn("failed to parse json", sl.Err(err))
		return
	}

	if note.Note == "" {
		http.Error(w, "Empty note", http.StatusBadRequest)
		h.log.Warn("invalid argument", sl.Err(errors.New("empty note text")))
		return
	}

	noteID, spellingErrors, err := h.notesService.CreateNote(r.Context(), note.Note, userID)
	if err != nil {
		http.Error(w, "Failed to add note: "+err.Error(), http.StatusBadRequest)
		h.log.Warn("failed to add note", sl.Err(err))
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"id":             noteID,
		"spellingErrors": spellingErrors,
	})
	if err != nil {
		http.Error(w, "Failed to marshal response json: "+err.Error(), http.StatusBadRequest)
		h.log.Warn("failed to marshal response json", sl.Err(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handler) getNotes(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 0 // for all notes
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 0 // for all notes
	}

	userID, ok := r.Context().Value(auth.CtxUserID).(int64)
	if !ok {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		h.log.Warn("invalid user id")
		return
	}

	notes, err := h.notesService.GetNotes(r.Context(), userID, page, limit)
	if err != nil {
		http.Error(w, "Failed to get notes: "+err.Error(), http.StatusBadRequest)
		h.log.Warn("failed to get notes", sl.Err(err))
		return
	}

	resp, err := json.Marshal(map[string][]models.Note{
		"notes": notes,
	})
	if err != nil {
		http.Error(w, "Failed to marshal response json: "+err.Error(), http.StatusBadRequest)
		h.log.Warn("failed to marshal response json", sl.Err(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

// func (h *Handler) getNotesPage(w http.ResponseWriter, r *http.Request) {
// 	pageString := r.URL.Query().Get("page")
// 	limitString := r.URL.Query().Get("limit")

// 	userID, ok := r.Context().Value(auth.CtxUserID).(int64)
// 	if !ok {
// 		http.Error(w, "Invalid user id", http.StatusBadRequest)
// 		h.log.Warn("invalid user id")
// 		return
// 	}

// 	page, err := strconv.Atoi(pageString)
// 	if err != nil {
// 		http.Error(w, "Failed to get query parameter page: "+err.Error(), http.StatusBadRequest)
// 		h.log.Warn("failed to get page", sl.Err(err))
// 		return
// 	}

// 	limit, err := strconv.Atoi(limitString)
// 	if err != nil {
// 		http.Error(w, "Failed to get query parameter limit: "+err.Error(), http.StatusBadRequest)
// 		h.log.Warn("failed to get limit", sl.Err(err))
// 		return
// 	}

// 	notes, err := h.notesService.GetNotesPage(r.Context(), userID, page, limit)
// 	if err != nil {
// 		http.Error(w, "Failed to get notes: "+err.Error(), http.StatusBadRequest)
// 		h.log.Warn("failed to get notes", sl.Err(err))
// 		return
// 	}

// 	resp, err := json.Marshal(map[string][]models.Note{
// 		"notes": notes,
// 	})
// 	if err != nil {
// 		http.Error(w, "Failed to marshal response json: "+err.Error(), http.StatusBadRequest)
// 		h.log.Warn("failed to marshal response json", sl.Err(err))
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Add("Content-Type", "application/json")
// 	w.Write(resp)
// }
