package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/piotmni/go-mini-templates/ent-atlas/ent"
	"github.com/piotmni/go-mini-templates/ent-atlas/ent/paste"
)

// PasteService handles paste business logic.
type PasteService struct {
	client *ent.Client
}

// NewPasteService creates a new paste service.
func NewPasteService(client *ent.Client) *PasteService {
	return &PasteService{client: client}
}

// CreatePasteInput holds data for creating a paste.
type CreatePasteInput struct {
	Title     string
	Content   string
	Language  string
	IsPublic  bool
	ExpiresIn *time.Duration
}

// PasteOutput represents a paste response.
type PasteOutput struct {
	ID        int        `json:"id"`
	Slug      string     `json:"slug"`
	Title     string     `json:"title,omitempty"`
	Content   string     `json:"content"`
	Language  string     `json:"language"`
	IsPublic  bool       `json:"is_public"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// Create creates a new paste.
func (s *PasteService) Create(ctx context.Context, input CreatePasteInput) (*PasteOutput, error) {
	slug, err := generateSlug(8)
	if err != nil {
		return nil, err
	}

	create := s.client.Paste.Create().
		SetSlug(slug).
		SetContent(input.Content).
		SetIsPublic(input.IsPublic)

	if input.Title != "" {
		create.SetTitle(input.Title)
	}

	if input.Language != "" {
		create.SetLanguage(input.Language)
	}

	if input.ExpiresIn != nil {
		expiresAt := time.Now().Add(*input.ExpiresIn)
		create.SetExpiresAt(expiresAt)
	}

	p, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	return toPasteOutput(p), nil
}

// GetBySlug retrieves a paste by its slug.
func (s *PasteService) GetBySlug(ctx context.Context, slug string) (*PasteOutput, error) {
	p, err := s.client.Paste.Query().
		Where(paste.Slug(slug)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	// Check if expired
	if p.ExpiresAt != nil && p.ExpiresAt.Before(time.Now()) {
		return nil, &ErrPasteExpired{Slug: slug}
	}

	return toPasteOutput(p), nil
}

// List retrieves recent public pastes.
func (s *PasteService) List(ctx context.Context, limit, offset int) ([]*PasteOutput, error) {
	pastes, err := s.client.Paste.Query().
		Where(paste.IsPublic(true)).
		Where(paste.Or(
			paste.ExpiresAtIsNil(),
			paste.ExpiresAtGT(time.Now()),
		)).
		Order(ent.Desc(paste.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*PasteOutput, len(pastes))
	for i, p := range pastes {
		result[i] = toPasteOutput(p)
	}

	return result, nil
}

// Delete removes a paste by slug.
func (s *PasteService) Delete(ctx context.Context, slug string) error {
	_, err := s.client.Paste.Delete().
		Where(paste.Slug(slug)).
		Exec(ctx)
	return err
}

func toPasteOutput(p *ent.Paste) *PasteOutput {
	return &PasteOutput{
		ID:        p.ID,
		Slug:      p.Slug,
		Title:     p.Title,
		Content:   p.Content,
		Language:  p.Language,
		IsPublic:  p.IsPublic,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
}

func generateSlug(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

// ErrPasteExpired indicates the paste has expired.
type ErrPasteExpired struct {
	Slug string
}

func (e *ErrPasteExpired) Error() string {
	return "paste has expired: " + e.Slug
}
