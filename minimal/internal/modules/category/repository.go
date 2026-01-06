package category

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("category not found")
	ErrAlreadyExists = errors.New("category already exists")
)

// Repository defines the interface for category persistence.
type Repository interface {
	Create(ctx context.Context, c Category) error
	GetByID(ctx context.Context, id ID) (Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, c Category) error
	Delete(ctx context.Context, id ID) error
}
