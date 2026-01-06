package note

import (
	"context"
	"time"

	"github.com/piotmni/go-mini-templates/minimal/internal/modules/category"
	"go.uber.org/zap"
)

// Service provides note business logic.
type Service struct {
	repo   Repository
	logger *zap.Logger
}

// NewService creates a new note service.
func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.Named("note.service"),
	}
}

// CreateInput contains data for creating a note.
type CreateInput struct {
	CategoryID category.ID
	Title      string
	Content    string
}

// Create creates a new note.
func (s *Service) Create(ctx context.Context, input CreateInput) (Note, error) {
	now := time.Now().UTC()
	n := Note{
		ID:         NewID(),
		CategoryID: input.CategoryID,
		Title:      input.Title,
		Content:    input.Content,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, n); err != nil {
		s.logger.Error("failed to create note", zap.Error(err))
		return Note{}, err
	}

	s.logger.Info("note created", zap.String("id", n.ID.String()))
	return n, nil
}

// GetByID retrieves a note by ID.
func (s *Service) GetByID(ctx context.Context, id ID) (Note, error) {
	n, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get note", zap.String("id", id.String()), zap.Error(err))
		return Note{}, err
	}
	return n, nil
}

// GetAll retrieves all notes.
func (s *Service) GetAll(ctx context.Context) ([]Note, error) {
	notes, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all notes", zap.Error(err))
		return nil, err
	}
	return notes, nil
}

// GetByCategory retrieves all notes in a category.
func (s *Service) GetByCategory(ctx context.Context, categoryID category.ID) ([]Note, error) {
	notes, err := s.repo.GetByCategory(ctx, categoryID)
	if err != nil {
		s.logger.Error("failed to get notes by category", zap.String("category_id", categoryID.String()), zap.Error(err))
		return nil, err
	}
	return notes, nil
}

// UpdateInput contains data for updating a note.
type UpdateInput struct {
	ID         ID
	CategoryID category.ID
	Title      string
	Content    string
}

// Update updates an existing note.
func (s *Service) Update(ctx context.Context, input UpdateInput) (Note, error) {
	n, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return Note{}, err
	}

	n.CategoryID = input.CategoryID
	n.Title = input.Title
	n.Content = input.Content
	n.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, n); err != nil {
		s.logger.Error("failed to update note", zap.String("id", input.ID.String()), zap.Error(err))
		return Note{}, err
	}

	s.logger.Info("note updated", zap.String("id", n.ID.String()))
	return n, nil
}

// Delete removes a note.
func (s *Service) Delete(ctx context.Context, id ID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete note", zap.String("id", id.String()), zap.Error(err))
		return err
	}
	s.logger.Info("note deleted", zap.String("id", id.String()))
	return nil
}
