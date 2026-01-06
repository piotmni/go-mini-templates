package note

import (
	"context"
	"errors"

	"github.com/piotmni/go-mini-templates/minimal/internal/modules/category"
)

var (
	ErrNotFound = errors.New("note not found")
)

// Repository defines the interface for note persistence.
type Repository interface {
	Create(ctx context.Context, n Note) error
	GetByID(ctx context.Context, id ID) (Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	GetByCategory(ctx context.Context, categoryID category.ID) ([]Note, error)
	Update(ctx context.Context, n Note) error
	Delete(ctx context.Context, id ID) error
}
