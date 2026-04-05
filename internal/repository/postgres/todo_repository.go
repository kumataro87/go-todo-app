package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-todo-app/internal/domain"

	"github.com/jmoiron/sqlx"
)

type todoRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *todoRepository {
	return &todoRepository{db: db}
}

type todoRow struct {
	domain.Todo
	TotalCount int `db:"total_count"`
}

func (r *todoRepository) GetAll(ctx context.Context, filter domain.ListFilter) (*domain.ListResult, error) {
	query := `
	SELECT id, title, description, completed, created_at, updated_at,
		COUNT(*) OVER() AS total_count
	FROM todos
	WHERE
		($1::boolean IS NULL OR completed = $1)
		AND ($2 = '' OR title ILIKE '%' || $2 || '%')
	ORDER BY created_at DESC
	LIMIT $3 OFFSET $4
	`

	rows := []todoRow{}
	err := r.db.SelectContext(ctx, &rows, query,
		filter.Completed,
		filter.Search,
		filter.Limit,
		filter.Offset(),
	)
	if err != nil {
		return nil, fmt.Errorf("get all todos: %w", err)
	}

	todos := make([]domain.Todo, 0, len(rows))
	totalCount := 0
	for _, row := range rows {
		todos = append(todos, row.Todo)
	}

	totalPages := 0
	if filter.Limit > 0 && totalCount > 0 {
		totalPages = (totalCount*filter.Limit - 1) / filter.Limit
	}

	return &domain.ListResult{
		Todos:      todos,
		TotalCount: totalCount,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *todoRepository) GetByID(ctx context.Context, id int) (*domain.Todo, error) {
	var todo domain.Todo
	err := r.db.GetContext(ctx, &todo,
		`SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1`,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError(fmt.Sprintf("todo not found: id=%d", id))
		}
		return nil, fmt.Errorf("get todo by id: %w", err)
	}
	return &todo, nil
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	var created domain.Todo
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO todos (title, description) VALUES ($1, $2)
		 RETURNING id, title, description, completed, created_at, updated_at`,
		todo.Title,
		todo.Description,
	).StructScan(&created)
	if err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}
	return &created, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	var updated domain.Todo
	err := r.db.QueryRowxContext(ctx,
		`UPDATE todos
		 SET title = $1, description = $2, completed = $3, updated_at = NOW()
		 WHERE id = $4
		 RETURNING id, title, description, completed, created_at, updated_at`,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.ID,
	).StructScan(&updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError(fmt.Sprintf("todo not found: id=%d", todo.ID))
		}
		return nil, fmt.Errorf("update todo: %w", err)
	}
	return &updated, nil
}

func (r *todoRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM todos WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return domain.NewNotFoundError(fmt.Sprintf("todo not found: id=%d", id))
	}
	return nil
}
