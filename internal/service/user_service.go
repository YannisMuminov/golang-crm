package service

import (
	"context"
	"fmt"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &UserService{repo: repo}
}

func (r *UserService) GetAll(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error) {
	filter.Sanitize()

	users, total, err := r.repo.GetAll(ctx, filter)

	if err != nil {
		return nil, total, err
	}

	return users, total, nil
}

func (r *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := r.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if user == nil {
		return nil, apperror.ErrNotFound
	}

	return user, nil
}

func (r *UserService) Update(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := r.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return nil, apperror.ErrNotFound
	}

	if req.FirstName != "" {
		req.FirstName = user.FirstName
	}

	if req.LastName != "" {
		req.LastName = user.LastName
	}

	if req.RoleID > 0 {
		user.RoleID = req.RoleID
	}

	return user, nil
}

func (r *UserService) Deactivate(ctx context.Context, id int64) error {
	user, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return apperror.ErrNotFound
	}

	if !user.IsActive {
		return nil
	}

	if err := r.repo.SetActive(ctx, id, false); err != nil {
		return fmt.Errorf("set active user: %w", err)
	}
	return nil
}

func (r *UserService) Activate(ctx context.Context, id int64) error {
	user, err := r.repo.GetByID(ctx, id)

	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return apperror.ErrNotFound
	}

	if !user.IsActive {
		return nil
	}

	if err := r.repo.SetActive(ctx, id, true); err != nil {
		return fmt.Errorf("set active user: %w", err)
	}

	return nil
}
