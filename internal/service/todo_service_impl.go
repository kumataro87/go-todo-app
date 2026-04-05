package service

import (
	"context"
	"log/slog"

	"go-todo-app/internal/domain"
	"go-todo-app/internal/repository"
)

type todoService struct {
	repo   repository.TodoRepository
	logger *slog.Logger
}

// New はTodoServiceの実装を返す
func New(repo repository.TodoRepository, logger *slog.Logger) TodoService {
	return &todoService{repo: repo, logger: logger}
}

func (s *todoService) GetAll(ctx context.Context, filter domain.ListFilter) (*domain.ListResult, error) {
	// デフォルト値の補完
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	result, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get all todos", slog.Any("error", err))
		return nil, domain.NewInternalError(err)
	}
	return result, nil
}

func (s *todoService) GetByID(ctx context.Context, id int) (*domain.Todo, error) {
	if id <= 0 {
		return nil, domain.NewBadRequestError("id must be a positive integer")
	}

	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err // リポジトリが返す AppError をそのまま伝播
	}
	return todo, nil
}

func (s *todoService) Create(ctx context.Context, input CreateInput) (*domain.Todo, error) {
	todo := &domain.Todo{
		Title:       input.Title,
		Description: input.Description,
	}

	created, err := s.repo.Create(ctx, todo)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to create todo", slog.Any("error", err))
		return nil, domain.NewInternalError(err)
	}

	s.logger.InfoContext(ctx, "todo created", slog.Int("id", created.ID))
	return created, nil
}

func (s *todoService) Update(ctx context.Context, id int, input UpdateInput) (*domain.Todo, error) {
	if id <= 0 {
		return nil, domain.NewBadRequestError("id must be a positive integer")
	}

	// 現在の値を取得（存在しなければ NotFound が返る）
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 送られてきたフィールドだけ上書き（部分更新）
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Completed != nil {
		existing.Completed = *input.Completed
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to update todo", slog.Int("id", id), slog.Any("error", err))
		return nil, err
	}

	s.logger.InfoContext(ctx, "todo updated", slog.Int("id", id))
	return updated, nil
}

func (s *todoService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.NewBadRequestError("id must be a positive integer")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "todo deleted", slog.Int("id", id))
	return nil
}
