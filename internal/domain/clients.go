package domain

import (
	"context"
	"time"
)

type Client struct {
	ID        int64     `db:"id"         json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name"  json:"last_name"`
	Email     string    `db:"email"      json:"email"`
	Phone     string    `db:"phone"      json:"phone"`
	Company   string    `db:"company"    json:"company"`
	Notes     string    `db:"notes"      json:"notes"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ClientCreateRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
	Notes     string `json:"notes"`
}

type ClientUpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
	Notes     string `json:"notes"`
}
type ClientFilter struct {
	Search string
	Page   int
	Limit  int
}

type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id int64) (*Client, error)
	GetByEmail(ctx context.Context, email string) (*Client, error)
	GetAll(ctx context.Context, filter ClientFilter) ([]Client, int, error)
	Update(ctx context.Context, client *Client) error
	Delete(ctx context.Context, id int64) error
}

type ClientService interface {
	Create(ctx context.Context, req *ClientCreateRequest, createdBy int64) (*Client, error)
	GetByID(ctx context.Context, id int64) (*Client, error)
	GetAll(ctx context.Context, filter ClientFilter) ([]Client, int, error)
	Update(ctx context.Context, id int64, req *ClientUpdateRequest) (*Client, error)
	Delete(ctx context.Context, id int64) error
}

func (f *ClientFilter) Offset() int {
	if f.Page <= 1 {
		return 0
	}
	return (f.Page - 1) * f.Limit
}

func (f *ClientFilter) Sanitize() {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
}
