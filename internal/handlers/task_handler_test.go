package handlers

import(
	"os"
	"net/http"
	"net/http/test"
	"testing"
	"bytes"
	"Task Manager API project/internal/storage"
)

func TestCreateTest(t *testing.T){
	store := storage.NewFileStore("test_tasks.json")
	defer os.Remove("test_tasks.json")

	handler := NewTaskHandler(store)

	body := []byte('{"title":"Learn Go", "completed":"false"}')
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	w := http.test.NewRecoder()

	handler.Tasks(w, req)

	if w.Code != http.StatusOk{
		t.Fatalf("expected 200, got %d", w.Code)
	}
}