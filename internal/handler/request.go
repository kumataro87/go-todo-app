package handler

import (
	"encoding/json"
	"fmt"
	"go-todo-app/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=1, max=255"`
	Description string `json:"description" validate:"max=1024"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1, max=255"`
	Description *string `json:"description" validate:"omitempty,max=1024"`
	Completed   *bool   `json:"completed"`
}

type ListQueryParams struct {
	Page      int    `validate:"min=1"`
	Limit     int    `validate:"min=1,max=100"`
	Completed *bool  `validate:"omitempty"`
	Search    string `validate:"max=255"`
}

func decodeJSON[T any](r *http.Request, dst *T) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return domain.NewBadRequestError("invalid JSON: " + err.Error())
	}
	if err := validate.Struct(dst); err != nil {
		return domain.NewBadRequestError(formatValidationErrors(err))
	}

	return nil
}

func parseListQuery(r *http.Request) (ListQueryParams, error) {
	q := r.URL.Query()

	params := ListQueryParams{
		Page:   1,
		Limit:  20,
		Search: q.Get("search"),
	}

	if v := q.Get("page"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return params, domain.NewBadRequestError("page must be a positive integer")
		}
		params.Page = n
	}

	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 100 {
			return params, domain.NewBadRequestError("limit must be between 1 and 100")
		}
		params.Limit = n
	}

	if v := q.Get("completed"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return params, domain.NewBadRequestError("completed must be true or false")
		}
		params.Completed = &b
	}

	return params, nil
}

func formatValidationErrors(err error) string {
	var errs validator.ValidationErrors
	if ve, ok := err.(validator.ValidationErrors); ok {
		errs = ve
	}
	if len(errs) == 0 {
		return err.Error()
	}
	msg := ""
	for i, fe := range errs {
		if i > 0 {
			msg += "; "
		}
		msg += fmt.Sprintf("field '%s' failed '%s' validation", fe.Field(), fe.Tag())
	}
	return msg
}
