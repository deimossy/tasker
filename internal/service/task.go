package service

import (
	"context"
	"fmt"
	"time"

	"github.com/deimossy/tasker/internal/entity"
	"github.com/deimossy/tasker/internal/repository"
)

type TaskService struct {
	repo   repository.TaskRepository
	logger chan string
}

func NewTaskService(repo repository.TaskRepository, logger chan string) *TaskService {
	go func() {
		for msg := range logger {
			fmt.Println(time.Now().Format(time.RFC3339), msg)
		}
	}()

	return &TaskService{
		repo:   repo,
		logger: logger,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	if !task.Status.IsValid() {
		return entity.Task{}, repository.ErrInvalidStatus
	}

	createdTask, err := s.repo.Create(ctx, task)
	if err != nil {
		return entity.Task{}, err
	}

	select {
	case s.logger <- fmt.Sprintf("Created task ID: %d", createdTask.ID):
	default:
	}

	return createdTask, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (entity.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}

	select {
	case s.logger <- fmt.Sprintf("Fetched task ID: %d", task.ID):
	default:
	}

	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, status *entity.Status) ([]entity.Task, error) {
	tasks, err := s.repo.List(ctx, status)
	if err != nil {
		return nil, err
	}

	select {
	case s.logger <- fmt.Sprintf("Listed tasks, count: %d", len(tasks)):
	default:
	}

	return tasks, nil
}
