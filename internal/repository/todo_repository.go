package repository

import (
	"context"
	"go-todo-app/internal/domain"
)

type TodoRepository interface {
	GetAll(ctx context.Context, filter domain.ListFilter) (*domain.ListResult, error)
	GetByID(ctx context.Context, id int) (*domain.Todo, error)
	Create(ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	Update(ctx context.Context, todo *domain.Todo) (*domain.Todo, error)
	Delete(ctx context.Context, id int) error
}
