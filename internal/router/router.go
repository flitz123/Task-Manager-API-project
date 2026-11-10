package router

import (
	"Task_Manager_API_project/internal/handlers"
	"Task_Manager_API_project/internal/middleware"
	"net/http"
	"os"
)

func Setup(handler *handlers.TaskHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", handler.Tasks)
	mux.HandleFunc("/tasks/", handler.TaskByID)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	if _, err := os.Stat("public"); err == nil {
		mux.Handle("/", http.FileServer(http.Dir("public")))
	}
	return middleware.Logging(middleware.CORS(middleware.Recovery(mux)))
}
