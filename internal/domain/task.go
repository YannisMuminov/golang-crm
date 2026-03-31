package domain

import (
	"context"
	"time"
)

type TaskStatus string

const (
	TaskStatusNew        TaskStatus = "new"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusNew, TaskStatusInProgress, TaskStatusDone:
		return true
	}
	return false
}

type Task struct {
	ID          int64      `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Status      TaskStatus `db:"status" json:"status"`
	DealID      int64      `db:"deal_id" json:"deal_id"`
	AssignedTo  int64      `db:"assigned_to" json:"assigned_to"`
	CreatedBy   int64      `db:"created_by" json:"created_by"`
	DueDate     *time.Time `db:"due_data" json:"due_data"` //Срок выполнение
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DealID      int64      `json:"deal_id" binding:"required"`
	AssignedTo  int64      `json:"assigned_to" binding:"required"`
	Status      TaskStatus `json:"status"`
	DueDate     *time.Time `json:"due_data"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	AssignedTo  int64      `json:"assigned_to"`
	DueDate     *time.Time `json:"due_data"`
}

type TaskFilter struct {
	DealID     int64
	AssignedTo int64
	Status     TaskStatus
	Page       int
	Limit      int
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id int64) (*Task, error)
	GetAll(ctx context.Context, filter TaskFilter) ([]Task, int, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int64) error
}

type TaskService interface {
	Create(ctx context.Context, req *CreateTaskRequest, createdBy int64) (*Task, error)
	GetByID(ctx context.Context, id int64) (*Task, error)
	GetAll(ctx context.Context, filter TaskFilter) ([]Task, int, error)
	Update(ctx context.Context, id int64, req *UpdateTaskRequest) (*Task, error)
	Delete(ctx context.Context, id int64) error
}

func (f *TaskFilter) Offset() int {
	if f.Page <= 1 {
		return 0
	}
	return (f.Page - 1) * f.Limit
}

func (f *TaskFilter) Sanitize() {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}

	if f.Page <= 0 {
		f.Page = 1
	}
}
