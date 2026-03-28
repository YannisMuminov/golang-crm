package domain

import "context"

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	RoleID    int64  `json:"role_id"`
}

type UserFilter struct {
	RoleID   int64
	IsActive *bool
	Page     int
	Limit    int
}

type UserRepository interface {
	GetAll(ctx context.Context, filter UserFilter) ([]User, int, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	SetActive(ctx context.Context, id int64, isActive *bool) error
}

type UserService interface {
	GetAll(ctx context.Context, filter UserFilter) ([]User, int, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, id int64, req *UpdateUserRequest) (*User, error)
	Deactivate(ctx context.Context, id int64) error
	Activate(ctx context.Context, id int64) error
}

func (u *UserFilter) Offset() int {

	if u.Page <= 1 {
		return 0
	}

	return (u.Page - 1) * u.Limit
}

func (u *UserFilter) Sanitize() {

	if u.Limit <= 0 || u.Limit >= 100 {
		u.Limit = 20
	}
	if u.Page <= 0 {
		u.Page = 1
	}
}
