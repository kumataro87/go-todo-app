package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"go-todo-app/internal/domain"
	"go-todo-app/internal/service"
)

type TodoHandler struct {
	svc    service.TodoService
	logger *slog.Logger
}

func New(svc service.TodoService, logger *slog.Logger) *TodoHandler {
	return &TodoHandler{svc: svc, logger: logger}
}

func (h *TodoHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/todos", h.handleGetAll)
	mux.HandleFunc("GET /api/v1/todos/{id}", h.handleGetByID)
	mux.HandleFunc("POST /api/v1/todos", h.handleCreate)
	mux.HandleFunc("PUT /api/v1/todos/{id}", h.handleUpdate)
	mux.HandleFunc("DELETE /api/v1/todos/{id}", h.handleDelete)
}

func (h *TodoHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	params, err := parseListQuery(r)
	if err != nil {
		writeError(w, r, err)
		return
	}

	filter := domain.ListFilter{
		Page:      params.Page,
		Limit:     params.Limit,
		Completed: params.Completed,
		Search:    params.Search,
	}

	result, err := h.svc.GetAll(r.Context(), filter)
	if err != nil {
		writeError(w, r, err)
		return
	}

	todos := result.Todos
	if todos == nil {
		todos = []domain.Todo{}
	}

	writeJSON(w, http.StatusOK, PageResponse[domain.Todo]{
		Data: todos,
		Meta: Meta{
			Page:       result.Page,
			Limit:      result.Limit,
			TotalCount: result.TotalCount,
			TotalPages: result.TotalPages,
		},
	})
}

func (h *TodoHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}

	todo, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeOK(w, todo)
}

func (h *TodoHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, err)
		return
	}

	todo, err := h.svc.Create(r.Context(), service.CreateInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeCreated(w, todo)
}

func (h *TodoHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}

	var req UpdateTodoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, r, err)
		return
	}

	todo, err := h.svc.Update(r.Context(), id, service.UpdateInput{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}
	writeOK(w, todo)
}

func (h *TodoHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, r, err)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, r, err)
		return
	}
	writeNoContent(w)
}

func parseID(r *http.Request) (int, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0, domain.NewBadRequestError("id must be a positive integer")
	}
	return id, nil
}
