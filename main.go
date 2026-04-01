package main

import (
	"go-todo-app/internal/handler"
	"go-todo-app/internal/middleware"
	"go-todo-app/internal/store"
	"log/slog"
	"net/http"
)

func main() {
	todoStore := store.NewMemoryStore()
	todoHandler := handler.NewTodoHandler(todoStore)

	mux := http.NewServeMux()
	todoHandler.RegisterRoutes(mux)

	var h http.Handler = mux
	h = middleware.Logger(slog.Default())(h)

	addr := ":8080"
	if err := http.ListenAndServe(addr, h); err != nil {
		slog.Error("server failed", "error", err)
	}
}
