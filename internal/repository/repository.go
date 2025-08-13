package repository

import (
	"context"

	"github.com/deimossy/tasker/internal/entity"
)

type Err string

var (
	ErrTaskNotFound = Err("task not found")
	ErrInvalidStatus = Err("invalid status")
)

func (e Err) Error() string {
	return string(e)
}

type TaskRepository interface {
	Create(ctx context.Context, task entity.Task) (entity.Task, error)
	GetByID(ctx context.Context, id int64) (entity.Task, error)
	List(ctx context.Context) ([]entity.Task, error)
}