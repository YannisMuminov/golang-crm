package domain

import (
	"context"
	"time"
)

type DealStatus string

const (
	DealStatusNew         DealStatus = "new"
	DealStatusNegotiation DealStatus = "negotiation"
	DealStatusWon         DealStatus = "won"
	DealStatusLost        DealStatus = "lost"
)

func (s DealStatus) IsValid() bool {
	switch s {
	case DealStatusNew, DealStatusNegotiation, DealStatusWon, DealStatusLost:
		return true
	}
	return false

}

func (s DealStatus) IsClosed() bool {
	return s == DealStatusWon || s == DealStatusLost
}

type Deal struct {
	ID         int64      `db:"id" json:"id"`
	Title      string     `db:"title" json:"title"`
	Amount     float64    `db:"amount" json:"amount"`
	Status     DealStatus `db:"status" json:"status"`
	ClientID   int64      `db:"client_id" json:"client_id"`
	AssignedTo int64      `db:"assigned_to" json:"assigned_to"`
	CreatedBy  int64      `db:"created_by" json:"created_by"`
	ClosedAt   *time.Time `db:"closed_at" json:"closed_at"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

type CreateDealRequest struct {
	Title      string     `json:"title" binding:"required"`
	Amount     float64    `json:"amount" binding:"required,min=0"`
	ClientID   int64      `json:"client_id" binding:"required"`
	AssignedTo int64      `json:"assigned_to" binding:"required"`
	Status     DealStatus `json:"status"`
}

type UpdateDealRequest struct {
	Title      string     `json:"title"`
	Amount     float64    `json:"amount" binding:"omitempty,min=0"`
	ClientID   int64      `json:"client_id"`
	AssignedTo int64      `json:"assigned_to"`
	Status     DealStatus `json:"status"`
}

type DealFilter struct {
	ClientID   int64
	AssignedTo int64
	Status     DealStatus
	Page       int
	Limit      int
}

type DealRepository interface {
	Create(ctx context.Context, deal *Deal) error
	GetByID(ctx context.Context, id int64) (*Deal, error)
	GetAll(ctx context.Context, filter DealFilter) ([]Deal, int, error)
	Update(ctx context.Context, deal *Deal) error
	Delete(ctx context.Context, id int64) error
}

type DealService interface {
	Create(ctx context.Context, req *CreateDealRequest, createdBy int64) (*Deal, error)
	GetByID(ctx context.Context, id int64) (*Deal, error)
	GetAll(ctx context.Context, filter DealFilter) ([]Deal, int, error)
	Update(ctx context.Context, id int64, req *UpdateDealRequest) (*Deal, error)
	Delete(ctx context.Context, id int64) error
}

func (f *DealFilter) Offset() int {
	if f.Page <= 1 {
		return 0
	}
	return (f.Page - 1) * f.Limit
}

func (f *DealFilter) Sanitize() {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
}
