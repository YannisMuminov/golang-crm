package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/YannisMuminov/internal/domain"
	"github.com/jmoiron/sqlx"
)

type dealRepository struct {
	db *sqlx.DB
}

func NewDealRepository(db *sqlx.DB) domain.DealRepository {
	return &dealRepository{db: db}
}

func (d *dealRepository) Create(ctx context.Context, deal *domain.Deal) error {
	query := `
		INSERT INTO deals (title, amount, status, client_id, assigned_to, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := d.db.QueryRowContext(
		ctx, query,
		deal.Title,
		deal.Amount,
		deal.Status,
		deal.ClientID,
		deal.AssignedTo,
		deal.CreatedBy,
	).Scan(&deal.ID, &deal.CreatedAt, &deal.UpdatedAt)

	if err != nil {
		return fmt.Errorf("create deal: %w", err)
	}

	return nil
}

func (d *dealRepository) GetByID(ctx context.Context, id int64) (*domain.Deal, error) {
	query := `
		SELECT id, title, amount, status, client_id, assigned_to, created_by, closed_at, created_at, updated_at
		FROM deals
		WHERE id = $1
	`

	var deal domain.Deal
	err := d.db.QueryRowContext(ctx, query, id).Scan(
		&deal.ID,
		&deal.Title,
		&deal.Amount,
		&deal.Status,
		&deal.ClientID,
		&deal.AssignedTo,
		&deal.CreatedBy,
		&deal.ClosedAt,
		&deal.CreatedAt,
		&deal.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return &deal, nil
}

func (d *dealRepository) GetAll(ctx context.Context, filter domain.DealFilter) ([]domain.Deal, int, error) {
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if filter.ClientID > 0 {
		conditions = append(conditions, fmt.Sprintf("client_id = $%d", argIdx))
		args = append(args, filter.ClientID)
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

	countQuery := "SELECT COUNT(*) FROM deals" + where

	var total int

	if err := d.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count deals: %w", err)
	}

	dataQuery := fmt.Sprintf(
		`
		SELECT id, title, amount, status, client_id, assigned_to, created_by, closed_at, created_at, updated_at
		FROM deals%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)

	args = append(args, filter.Limit, filter.Offset())

	var deals []domain.Deal

	if err := d.db.SelectContext(ctx, &deals, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("get all deals: %w", err)
	}
	return deals, total, nil

}

func (d *dealRepository) Update(ctx context.Context, deal *domain.Deal) error {
	query := `
		UPDATE deals
		SET title = $1,
			amount = $2,
			status = $3,
			assigned_to = $4,
			closed_at = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := d.db.QueryRowContext(
		ctx, query,
		deal.Title,
		deal.Amount,
		deal.Status,
		deal.AssignedTo,
		deal.ClosedAt,
		deal.ID,
	).Scan(&deal.UpdatedAt)

	if err != nil {
		return fmt.Errorf("update deal: %w", err)
	}

	return nil
}

func (d *dealRepository) Delete(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM deals WHERE id = $1`, id)

	if err != nil {
		return fmt.Errorf("delete deal: %w", err)
	}

	return nil
}
