package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go-todo-app/internal/domain"
)

type Response[T any] struct {
	Data  T      `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type PageResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	appErr := domain.AsAppError(err)

	if appErr.Code >= 500 {
		slog.ErrorContext(r.Context(), "internal error",
			slog.Any("error", appErr.Err),
		)
	}

	writeJSON(w, appErr.Code, Response[any]{Error: appErr.Message})
}

func writeOK(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusOK, Response[any]{Data: data})
}

func writeCreated(w http.ResponseWriter, data any) {
	writeJSON(w, http.StatusCreated, Response[any]{Data: data})
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}
