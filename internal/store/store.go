package store

import "go-todo-app/internal/model"

type TodoStore interface {
	GetAll() ([]model.Todo, error)
	GetById(id int) (model.Todo, error)
	Create(todo model.Todo) (model.Todo, error)
	Update(id int, todo model.Todo) (model.Todo, error)
	Delete(id int) error
}
