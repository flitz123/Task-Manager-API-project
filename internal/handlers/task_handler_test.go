package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Task_Manager_API_project/internal/storage"
)

func TestCreateTask(t *testing.T) {
	store := storage.NewFileStore("test_tasks.json")
	defer os.Remove("test_tasks.json")
	handler := NewTaskHandler(store)

	body := []byte(`{"title":"Learn Go","completed":false}`)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.Tasks(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}
