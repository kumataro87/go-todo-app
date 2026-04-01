package store

import (
	"fmt"
	"go-todo-app/internal/model"
	"sync"
	"time"
)

type MemoryStore struct {
	mu     sync.RWMutex
	todos  map[int]model.Todo
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		todos:  make(map[int]model.Todo),
		nextID: 1,
	}
}

func (s *MemoryStore) GetAll() ([]model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]model.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s *MemoryStore) GetById(id int) (model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("todo not found: id=%d", id)
	}
	return todo, nil
}

func (s *MemoryStore) Create(todo model.Todo) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo.ID = s.nextID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	s.todos[todo.ID] = todo
	s.nextID++
	return todo, nil
}

func (s *MemoryStore) Update(id int, todo model.Todo) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.todos[id]
	if !ok {
		return model.Todo{}, fmt.Errorf("todo not found: id=%d", id)
	}

	existing.Title = todo.Title
	existing.Completed = todo.Completed
	existing.UpdatedAt = time.Now()
	s.todos[id] = existing

	return existing, nil
}

func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("todo not found: id=%d", id)
	}
	delete(s.todos, id)
	return nil
}
