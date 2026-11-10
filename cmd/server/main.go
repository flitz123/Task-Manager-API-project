package main

import (
	"Task_Manager_API_project/internal/handlers"
	"Task_Manager_API_project/internal/router"
	"Task_Manager_API_project/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	dataFile := os.Getenv("TASK_DATA_FILE")
	if dataFile == "" {
		dataFile = "data/tasks.json"
	}
	store := storage.NewFileStore(dataFile)
	handler := handlers.NewTaskHandler(store)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router.Setup(handler),
	}

	log.Println("Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
