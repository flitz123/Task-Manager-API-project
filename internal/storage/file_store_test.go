package storage

import (
	"os"
	"testing"

	"Task_Manager_API_project/internal/models"
)

func TestAddAndGetTask(t *testing.T) {
	file := "test_tasks.json"
	defer os.Remove(file)
	store := NewFileStore(file)

	created := store.Add(models.Task{Title: "Test Task", Completed: false})
	if created.ID != 1 {
		t.Fatalf("expected ID 1, got %d", created.ID)
	}
	fetched, ok := store.GetByID(1)
	if !ok {
		t.Fatal("task not found")
	}
	if fetched.Title != "Test Task" {
		t.Errorf("expected 'Test Task', got %s", fetched.Title)
	}
}
