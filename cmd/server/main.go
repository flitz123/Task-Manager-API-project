package main

import (
	"Task Manager API project/internal/handlers"
	"Task Manager API project/internal/routers"
	"Task Manager API project/internal/storage"
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
	log.Ftal(server.ListenAndServe())
}
