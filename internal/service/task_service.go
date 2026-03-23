package service

import (
	"context"
	"fmt"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
)

type taskService struct {
	repoTask domain.TaskRepository
	repoDeal domain.DealRepository
}

func NewTaskService(repoTask domain.TaskRepository, repoDeal domain.DealRepository) domain.TaskService {
	return &taskService{repoTask: repoTask, repoDeal: repoDeal}
}

// Create implements [domain.TaskService].
func (s *taskService) Create(ctx context.Context, req *domain.CreateTaskRequest, createdBy int64) (*domain.Task, error) {
	deal, err := s.repoDeal.GetByID(ctx, req.DealID)

	if err != nil {
		return nil, fmt.Errorf("create task check deal: %w", err)
	}

	if deal == nil {
		return nil, apperror.ErrNotFound
	}

	status := req.Status

	if status == "" {
		status = domain.TaskStatusNew
	}
	if !status.IsValid() {
		return nil, apperror.ErrBadRequest
	}

	task := &domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DealID:      req.DealID,
		AssignedTo:  req.AssignedTo,
		CreatedBy:   createdBy,
		DueDate:     req.DueDate,
	}

	if err := s.repoTask.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

// GetByID implements [domain.TaskService].
func (s *taskService) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	task, err := s.repoTask.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	if task == nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return task, nil
}

// GetAll implements [domain.TaskService].
func (s *taskService) GetAll(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, int, error) {
	filter.Sanitize()

	if filter.Status != "" && !filter.Status.IsValid() {
		return nil, 0, apperror.ErrBadRequest
	}

	tasks, total, err := s.repoTask.GetAll(ctx, filter)

	if err != nil {
		return nil, 0, fmt.Errorf("get all tasks: %w", err)
	}

	return tasks, total, nil
}

// Update implements [domain.TaskService].
func (s *taskService) Update(ctx context.Context, id int64, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.repoTask.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("update task get: %w", err)
	}

	if task == nil {
		return nil, apperror.ErrNotFound
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.AssignedTo > 0 {
		task.AssignedTo = req.AssignedTo
	}

	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if req.Status != "" {
		if !req.Status.IsValid() {
			return nil, apperror.ErrBadRequest
		}

		task.Status = req.Status
	}

	if err := s.repoTask.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("update task: %w", err)
	}

	return task, err
}

// Delete implements [domain.TaskService].
func (s *taskService) Delete(ctx context.Context, id int64) error {
	task, err := s.repoTask.GetByID(ctx, id)

	if err != nil {
		return fmt.Errorf("delete get task: %w", err)
	}

	if task == nil {
		return apperror.ErrBadRequest
	}

	if err := s.repoTask.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
