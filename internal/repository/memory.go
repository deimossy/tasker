package repository

import (
	"context"
	"sync"
	"time"

	"github.com/deimossy/tasker/internal/entity"
)

type InMemoryTaskRepository struct {
	mu     sync.RWMutex
	data   map[int64]entity.Task
	nextID int64
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		data:   make(map[int64]entity.Task),
		nextID: 1,
	}
}

func (r *InMemoryTaskRepository) Create(ctx context.Context, task entity.Task) (entity.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-ctx.Done():
		return entity.Task{}, ctx.Err()
	default:
	}

	now := time.Now().UTC()
	task.ID = r.nextID
	task.CreatedAt = now
	task.UpdatedAt = now

	r.data[task.ID] = task
	r.nextID++

	return task, nil
}

func (r *InMemoryTaskRepository) GetByID(ctx context.Context, id int64) (entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return entity.Task{}, ctx.Err()
	default:
	}

	task, ok := r.data[id]
	if !ok {
		return entity.Task{}, ErrTaskNotFound
	}

	return task, nil
}

func (r *InMemoryTaskRepository) List(ctx context.Context, status *entity.Status) ([]entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	out := make([]entity.Task, 0, len(r.data))
	for _, task := range r.data {
		if status != nil && task.Status != *status {
			continue
		}
		out = append(out, task)
	}

	return out, nil
}
