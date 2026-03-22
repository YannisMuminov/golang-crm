package service

import (
	"context"
	"fmt"
	"time"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
)

type dealService struct {
	repo       domain.DealRepository
	clientRepo domain.ClientRepository
}

func NewDealService(repo domain.DealRepository, clientRepo domain.ClientRepository) domain.DealService {
	return &dealService{repo: repo, clientRepo: clientRepo}
}

func (d *dealService) Create(ctx context.Context, req *domain.CreateDealRequest, createdBy int64) (*domain.Deal, error) {
	client, err := d.clientRepo.GetByID(ctx, req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("create deal check client: %w", err)
	}

	if client == nil {
		return nil, apperror.ErrNotFound
	}

	status := req.Status

	if status == "" {
		status = domain.DealStatusNew
	}

	if !status.IsValid() {
		return nil, apperror.ErrBadRequest
	}

	deal := &domain.Deal{
		Title:      req.Title,
		Amount:     req.Amount,
		Status:     req.Status,
		ClientID:   req.ClientID,
		AssignedTo: req.AssignedTo,
		CreatedBy:  createdBy,
	}

	if status.IsClosed() {
		now := time.Now()
		deal.ClosedAt = &now
	}
	if err := d.repo.Create(ctx, deal); err != nil {
		return nil, fmt.Errorf("create deal: %w", err)
	}

	return deal, nil
}

func (d *dealService) GetByID(ctx context.Context, id int64) (*domain.Deal, error) {
	deal, err := d.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get deal by id: %w", err)
	}

	if deal == nil {
		return nil, apperror.ErrNotFound
	}

	return deal, nil
}

// GetAll implements [domain.DealService].
func (d *dealService) GetAll(ctx context.Context, filter domain.DealFilter) ([]domain.Deal, int, error) {
	filter.Sanitize()

	if filter.Status != "" && !filter.Status.IsValid() {
		return nil, 0, apperror.ErrBadRequest
	}
	deals, total, err := d.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("get all deals: %w", err)
	}
	return deals, total, nil
}

func (d *dealService) Update(ctx context.Context, id int64, req *domain.UpdateDealRequest) (*domain.Deal, error) {
	deal, err := d.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("get deal by id: %w", err)
	}

	if deal == nil {
		return nil, apperror.ErrNotFound
	}

	if req.Title != "" {
		deal.Title = req.Title
	}
	if req.Amount > 0 {
		deal.Amount = req.Amount
	}

	if req.AssignedTo > 0 {
		deal.AssignedTo = req.AssignedTo
	}

	if req.Status != "" {
		if !req.Status.IsValid() {
			return nil, apperror.ErrBadRequest
		}
		prevStatus := deal.Status
		deal.Status = req.Status

		if req.Status.IsClosed() && !prevStatus.IsClosed() {
			now := time.Now()
			deal.ClosedAt = &now
		}

		if !req.Status.IsClosed() && prevStatus.IsClosed() {
			deal.ClosedAt = nil
		}
	}

	if err := d.repo.Update(ctx, deal); err != nil {
		fmt.Println("UPDATE ERROR:", err)
		return nil, fmt.Errorf("update deal: %w", err)
	}

	return deal, nil

}

// Delete implements [domain.DealService].
func (d *dealService) Delete(ctx context.Context, id int64) error {
	deal, err := d.repo.GetByID(ctx, id)

	if err != nil {
		return fmt.Errorf("delete deal get: %w", err)
	}

	if deal == nil {
		return apperror.ErrNotFound
	}

	if err := d.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete deal: %w", err)
	}

	return nil
}
