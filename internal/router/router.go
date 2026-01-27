package router

import (
	"Task Manager API project/internal/handlers"
	"Task Manager API project/internal/middleware"
	"net/http"
)

func Setup(handler *handlers.TaskHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", handler.Tasks)
	mux.HandleFunc("/tasks/", handler.TaskByID)

	return middleware.Logging(middleware.Recovery(mux))
}
