package category

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// Service provides category business logic.
type Service struct {
	repo   Repository
	logger *zap.Logger
}

// NewService creates a new category service.
func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger.Named("category.service"),
	}
}

// CreateInput contains data for creating a category.
type CreateInput struct {
	Name string
}

// Create creates a new category.
func (s *Service) Create(ctx context.Context, input CreateInput) (Category, error) {
	now := time.Now().UTC()
	c := Category{
		ID:        NewID(),
		Name:      input.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		s.logger.Error("failed to create category", zap.Error(err))
		return Category{}, err
	}

	s.logger.Info("category created", zap.String("id", c.ID.String()))
	return c, nil
}

// GetByID retrieves a category by ID.
func (s *Service) GetByID(ctx context.Context, id ID) (Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get category", zap.String("id", id.String()), zap.Error(err))
		return Category{}, err
	}
	return c, nil
}

// GetAll retrieves all categories.
func (s *Service) GetAll(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all categories", zap.Error(err))
		return nil, err
	}
	return categories, nil
}

// UpdateInput contains data for updating a category.
type UpdateInput struct {
	ID   ID
	Name string
}

// Update updates an existing category.
func (s *Service) Update(ctx context.Context, input UpdateInput) (Category, error) {
	c, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return Category{}, err
	}

	c.Name = input.Name
	c.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, c); err != nil {
		s.logger.Error("failed to update category", zap.String("id", input.ID.String()), zap.Error(err))
		return Category{}, err
	}

	s.logger.Info("category updated", zap.String("id", c.ID.String()))
	return c, nil
}

// Delete removes a category.
func (s *Service) Delete(ctx context.Context, id ID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete category", zap.String("id", id.String()), zap.Error(err))
		return err
	}
	s.logger.Info("category deleted", zap.String("id", id.String()))
	return nil
}
