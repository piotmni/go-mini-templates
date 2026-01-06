package note

import (
	"time"

	"github.com/google/uuid"
	"github.com/piotmni/go-mini-templates/minimal/internal/modules/category"
)

// ID represents a note identifier.
type ID = uuid.UUID

// Note represents a note domain model.
type Note struct {
	ID         ID
	CategoryID category.ID
	Title      string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewID generates a new note ID.
func NewID() ID {
	return uuid.New()
}

// ParseID parses a string into a note ID.
func ParseID(s string) (ID, error) {
	return uuid.Parse(s)
}
