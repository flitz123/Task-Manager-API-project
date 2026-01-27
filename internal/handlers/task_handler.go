package handlers

import(
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"Task Manager API project/internal/models"
	"Task Manager API project/internal/storage"
)

type TaskHandler struct{
	store *storage.FileStore
}

func NewTaskHandler(store *storage.FileStore) *TaskHandler{
	return %TaskHandler{store: store}
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request){
	switch r.Method{
	case http.MethodGet:
		json.NewEncoder(w).Encode(h.Store.GetAll())

	case http.MethodPost:
		var tasks models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil{
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		created := h.store.Add(task)
		json.NewEncoder(w).Encode(created)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func(h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request){
	idStr := strings.TrimPrefix(r.URL.Path,"/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil{
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method{
	case http.MethodGet:
		task, ok:= h.store.GetByID(task)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(task)

	case http.MethodPut:
		var tasks models.Task
		_= json.NewDecoder(r.Body).Decode(&task) 
		updated, ok := h,store.update(id, task)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(updated)

	case http.MethodDelete:
		if h.store.Delete(id){
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:

		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

