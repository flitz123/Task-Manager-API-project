package handlers

import (
	"Task_Manager_API_project/internal/models"
	"Task_Manager_API_project/internal/storage"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	store *storage.FileStore
}

func NewTaskHandler(store *storage.FileStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func decodeTask(w http.ResponseWriter, r *http.Request) (models.Task, bool) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		return task, false
	}
	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Title is required"})
		return task, false
	}
	switch task.Priority {
	case "", "low", "medium", "high":
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Priority must be low, medium, or high"})
		return task, false
	}
	return task, true
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, h.store.GetAll())
	case http.MethodPost:
		task, ok := decodeTask(w, r)
		if !ok {
			return
		}
		created := h.store.Add(task)
		writeJSON(w, http.StatusCreated, created)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		task, ok := h.store.GetByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, task)
	case http.MethodPut, http.MethodPatch:
		task, ok := decodeTask(w, r)
		if !ok {
			return
		}
		updated, ok := h.store.Update(id, task)
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		if !h.store.Delete(id) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
