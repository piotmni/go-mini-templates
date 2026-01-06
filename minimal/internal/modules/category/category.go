package category

import (
	"time"

	"github.com/google/uuid"
)

// ID represents a category identifier.
type ID = uuid.UUID

// Category represents a note category domain model.
type Category struct {
	ID        ID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewID generates a new category ID.
func NewID() ID {
	return uuid.New()
}

// ParseID parses a string into a category ID.
func ParseID(s string) (ID, error) {
	return uuid.Parse(s)
}
