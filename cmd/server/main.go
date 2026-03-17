package main

import (
	"Task_Manager_API_project/internal/handlers"
	"Task_Manager_API_project/internal/router"
	"Task_Manager_API_project/internal/storage"
	"log"
	"net/http"
)

func main() {
	store := storage.NewFileStore("data/tasks.json")
	handler := handlers.NewTaskHandler(store)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router.Setup(handler),
	}

	log.Println("Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
