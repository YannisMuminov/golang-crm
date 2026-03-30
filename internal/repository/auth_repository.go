package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/YannisMuminov/internal/domain"
	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) domain.AuthRepository {
	return &authRepository{db: db}
}

func (a authRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (email, password_hash, first_name, last_name, role_id) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at`

	err := a.db.QueryRowContext(
		ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.RoleID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil

}

func (a authRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
	SELECT id, email, password_hash, first_name, last_name, role_id, is_active, created_at, updated_at
FROM users
	WHERE email = $1`

	var user domain.User

	err := a.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (a authRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role_id, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user domain.User

	err := a.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.RoleID,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (a authRepository) GetRolePermissions(ctx context.Context, roleID int64) (*domain.Role, []domain.Permission, error) {
	roleQuery := `
	SELECT id, name, description, created_at, updated_at
	FROM roles
	WHERE id = $1`

	var role domain.Role
	err := a.db.QueryRowContext(ctx, roleQuery, roleID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("get role permissions: %w", err)
	}

	permQuery := `
	SELECT p.id, p.name, p.description, p.created_at
	FROM permissions p
	JOIN role_permissions rp ON rp.permission_id = p.id
	WHERE rp.role_id = $1`

	var permissions []domain.Permission
	if err := a.db.SelectContext(ctx, &permissions, permQuery, role.ID); err != nil {
		return nil, nil, fmt.Errorf("get permissions: %w", err)
	}
	return &role, permissions, nil
}
