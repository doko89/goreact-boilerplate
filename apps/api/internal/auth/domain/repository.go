package domain

import "context"

// Repository persists users. Implemented in infrastructure.
type Repository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, u *User) error
	Count(ctx context.Context) (int64, error)
}
