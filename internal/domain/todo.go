package domain

import "time"

type Todo struct {
	ID          int       `db:"id"          json:"id"`
	Title       string    `db:"title"       json:"title"`
	Description string    `db:"description" json:"description"`
	Completed   bool      `db:"completed"   json:"completed"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
}

type ListFilter struct {
	Page      int
	Limit     int
	Completed *bool
	Search    string
}

type ListResult struct {
	Todos      []Todo `json:"todos"`
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

func (f ListFilter) Offset() int {
	if f.Page <= 0 {
		return 0
	}
	return (f.Page - 1) * f.Limit
}
