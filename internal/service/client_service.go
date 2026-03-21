package service

import (
	"context"
	"fmt"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
)

type clientService struct {
	repo domain.ClientRepository
}

func NewClientService(repo domain.ClientRepository) domain.ClientService {
	return &clientService{repo: repo}
}

func (c *clientService) Create(ctx context.Context, req *domain.ClientCreateRequest, createdBy int64) (*domain.Client, error) {
	existing, err := c.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	if existing != nil {
		return nil, apperror.ErrEmailTaken
	}

	client := &domain.Client{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Company:   req.Company,
		Notes:     req.Notes,
		CreatedBy: createdBy,
	}
	if err := c.repo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return client, nil
}

func (c *clientService) GetByID(ctx context.Context, id int64) (*domain.Client, error) {
	client, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}
	if client == nil {
		return nil, apperror.ErrNotFound
	}

	return client, nil
}

func (c *clientService) GetAll(ctx context.Context, filter domain.ClientFilter) ([]domain.Client, int, error) {
	filter.Sanitize()

	clients, total, err := c.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, total, err
	}
	return clients, total, nil
}

func (c *clientService) Update(ctx context.Context, id int64, req *domain.ClientUpdateRequest) (*domain.Client, error) {
	client, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}
	if client == nil {
		return nil, apperror.ErrNotFound
	}

	if req.FirstName != "" {
		client.FirstName = req.FirstName
	}
	if req.LastName != "" {
		client.LastName = req.LastName
	}
	if req.Phone != "" {
		client.Phone = req.Phone
	}
	if req.Company != "" {
		client.Company = req.Company
	}
	if req.Notes != "" {
		client.Notes = req.Notes
	}

	if req.Email != "" && req.Email != client.Email {
		existing, err := c.repo.GetByEmail(ctx, client.Email)
		if err != nil {
			return nil, fmt.Errorf("update client: %w", err)
		}
		if existing != nil {
			return nil, apperror.ErrEmailTaken
		}
		client.Email = req.Email
	}
	if err := c.repo.Update(ctx, client); err != nil {
		return nil, fmt.Errorf("update client: %w", err)
	}
	return client, nil
}

func (c *clientService) Delete(ctx context.Context, id int64) error {
	client, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}
	if client == nil {
		return apperror.ErrNotFound
	}
	if err := c.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete client: %w", err)
	}
	return nil
}
