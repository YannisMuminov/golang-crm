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

type clientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) domain.ClientRepository {
	if db == nil {
		panic(errors.New("db is nil"))
	}
	return &clientRepository{db: db}
}

func (c *clientRepository) Create(ctx context.Context, client *domain.Client) error {
	query := `
	INSERT INTO clients (first_name, last_name, email, phone, company, notes, created_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, updated_at`

	err := c.db.QueryRowContext(
		ctx, query,
		client.FirstName,
		client.LastName,
		client.Email,
		client.Phone,
		client.Company,
		client.Notes,
		client.CreatedBy,
	).Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	return nil
}

func (c *clientRepository) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	query := `
	SELECT id, first_name, last_name, email, phone, company, notes, created_by, created_at, updated_at
	FROM clients
	WHERE id = $1`

	var client domain.Client

	err := c.db.QueryRowContext(
		ctx, query, id,
	).Scan(
		&client.ID,
		&client.FirstName,
		&client.LastName,
		&client.Email,
		&client.Phone,
		&client.Company,
		&client.Notes,
		&client.CreatedBy,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get client: %w", err)
	}
	return &client, nil
}

func (c *clientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	query := `
	SELECT id, first_name, last_name, email, phone, company, notes, created_by, created_at, updated_at
	FROM clients
	WHERE email = $1`

	var client domain.Client

	err := c.db.QueryRowContext(
		ctx, query, email,
	).Scan(
		&client.ID,
		&client.FirstName,
		&client.LastName,
		&client.Email,
		&client.Phone,
		&client.Company,
		&client.Notes,
		&client.CreatedBy,
		&client.CreatedAt,
		&client.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get client by email: %w", err)
	}
	return &client, nil
}

func (c *clientRepository) GetAll(ctx context.Context, filter domain.ClientFilter) ([]domain.Client, int, error) {
	var args []interface{}
	argsIdx := 1
	var where strings.Builder

	if filter.Search != "" {
		where.WriteString(fmt.Sprintf(
			" WHERE (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d OR company ILIKE $%d)",
			argsIdx, argsIdx+1, argsIdx+2, argsIdx+3,
		))
		search := "%" + filter.Search + "%"
		args = append(args, search, search, search, search)
		argsIdx += 4
	}

	countQuery := "SELECT COUNT(*) FROM clients" + where.String()
	var total int

	if err := c.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("get all clients: %w", err)
	}

	dataQuery := fmt.Sprintf(`
	SELECT id, first_name, last_name, email, phone, company, notes, created_by, created_at, updated_at
	FROM clients%s
	ORDER BY created_at DESC
	LIMIT $%d OFFSET $%d`,
		where.String(), argsIdx, argsIdx+1,
	)
	args = append(args, filter.Limit, filter.Offset())
	var clients []domain.Client

	if err := c.db.SelectContext(ctx, &clients, dataQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("get all clients: %w", err)
	}

	return clients, total, nil
}

func (c *clientRepository) Update(ctx context.Context, client *domain.Client) error {
	query := `UPDATE clients
	SET first_name = $1,
    last_name = $2,
    email = $3,
    phone = $4,
    company = $5,
    notes = $6,
    updated_at = NOW()
    WHERE id = $7
    RETURNING updated_at`

	if err := c.db.QueryRowContext(
		ctx, query,
		&client.FirstName,
		&client.LastName,
		&client.Email,
		&client.Phone,
		&client.Company,
		&client.Notes,
		&client.UpdatedAt,
	).Scan(&client.UpdatedAt); err != nil {
		return fmt.Errorf("update client: %w", err)
	}
	return nil
}

func (c *clientRepository) Delete(ctx context.Context, id int64) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM clients WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete client: %w", err)
	}
	return nil
}
