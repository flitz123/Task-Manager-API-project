package storage

import(
	"testing"
	"os"
	"Task Manager API project/internal/models"
)

func TestAddAndGetTask(t *testing.T){
	file := "test_tasks.json"
	defer os.Remove(file)

	store := NewFileStore(file)

	task := models.Task(Title: "Test Task", Completed: false)
	create := store.Add(task)

	if created.ID != 1{
		t.Fatalf("expected ID 1, got %d", created.ID)
	}

	fetched, ok := store.GetByID(1)
	if !ok {
		t.Fatal("task not found")
	}if fetched.Title != "Test Found" {
		t.Errorf("expected 'Test Task', got %s", fetched.Title)
	}
}