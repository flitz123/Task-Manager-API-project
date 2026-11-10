package api

import (
	"net/http"
	"os"

	"Task_Manager_API_project/internal/handlers"
	"Task_Manager_API_project/internal/middleware"
	"Task_Manager_API_project/internal/storage"
)

var handler http.Handler

func init() {
	dataFile := os.Getenv("TASK_DATA_FILE")
	if dataFile == "" {
		dataFile = "/tmp/task-manager/tasks.json"
	}
	taskHandler := handlers.NewTaskHandler(storage.NewFileStore(dataFile))
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", taskHandler.Tasks)
	mux.HandleFunc("/tasks/", taskHandler.TaskByID)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	handler = middleware.CORS(middleware.Logging(middleware.Recovery(mux)))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r)
}
