package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/YannisMuminov/internal/domain"
	"github.com/jmoiron/sqlx"
)

type taskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (title, description, status, deal_id, assigned_to, created_by, due_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.DealID,
		task.AssignedTo,
		task.CreatedBy,
		task.DueDate,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	return nil
}

func (r *taskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, deal_id, assigned_to, created_by, due_data, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	var task domain.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.DealID,
		&task.AssignedTo,
		&task.CreatedBy,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}

	return &task, nil
}

// GetAll implements [domain.TaskRepository].
func (r *taskRepository) GetAll(ctx context.Context, filter domain.TaskFilter) ([]domain.Task, int, error) {
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if filter.DealID > 0 {
		conditions = append(conditions, fmt.Sprintf("deal_id = $%d", argIdx))
		args = append(args, filter.DealID)
		argIdx++
	}

	if filter.AssignedTo > 0 {
		conditions = append(conditions, fmt.Sprintf("assigned_to = $%d", argIdx))
		args = append(args, filter.AssignedTo)
		argIdx++
	}

	if filter.Status != "" && filter.Status.IsValid() {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}

	where := ""

	if len(conditions) > 0 {
		where = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM tasks" + where
	var total int

	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count tasks: %d", err)
	}

	dataQuery := fmt.Sprintf(
		`
		SELECT id, title, description, status, deal_id, assigned_to, created_by, due_data, created_at, updated_at
		FROM tasks%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	args = append(args, filter.Limit, filter.Offset())

	var tasks []domain.Task

	if err := r.db.SelectContext(ctx, &tasks, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("get all: %w", err)
	}
	return tasks, total, nil
}

// Update implements [domain.TaskRepository].
func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			assigned_to = $4,
			due_data = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.AssignedTo,
		task.DueDate,
		task.ID,
	).Scan(&task.UpdatedAt)

	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	return nil
}

// Delete implements [domain.TaskRepository].
func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	return nil
}
