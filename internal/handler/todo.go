package handler

import (
	"encoding/json"
	"go-todo-app/internal/model"
	"go-todo-app/internal/store"
	"net/http"
	"strconv"
	"strings"
)

type TodoHandler struct {
	store store.TodoStore
}

func NewTodoHandler(s store.TodoStore) *TodoHandler {
	return &TodoHandler{store: s}
}

func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /todos", h.handleGetAll)
	mux.HandleFunc("GET /todo/{id}", h.handleGetById)
	mux.HandleFunc("POST /todo", h.handleCreate)
	mux.HandleFunc("PUT /todo/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /todo/{id}", h.handleDelete)
}

func (h *TodoHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	todos, err := h.store.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, todos)
}

func (h *TodoHandler) handleGetById(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	todo, err := h.store.GetById(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJson(w, http.StatusOK, todo)
}

func (h *TodoHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo := model.Todo{Title: req.Title}
	todo, err := h.store.Create(todo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusCreated, todo)
}

func (h *TodoHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	existing, err := h.store.GetById(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var req model.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}

	if req.Completed != nil {
		existing.Completed = *req.Completed
	}

	updated, err := h.store.Update(id, existing)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, updated)
}

func (h *TodoHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	existing, err := h.store.GetById(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, existing)
}

func parseID(r *http.Request) (int, error) {
	// 文字列を整数に変換して返す
	return strconv.Atoi(r.PathValue("id"))
}

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJson(w, status, map[string]string{"error": message})
}
