package service

import (
	"context"

	"go-todo-app/internal/domain"
)

type TodoService interface {
	GetAll(ctx context.Context, filter domain.ListFilter) (*domain.ListResult, error)
	GetByID(ctx context.Context, id int) (*domain.Todo, error)
	Create(ctx context.Context, input CreateInput) (*domain.Todo, error)
	Update(ctx context.Context, id int, input UpdateInput) (*domain.Todo, error)
	Delete(ctx context.Context, id int) error
}

type CreateInput struct {
	Title       string
	Description string
}

type UpdateInput struct {
	Title       *string
	Description *string
	Completed   *bool
}
