package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/piotmni/go-mini-templates/minimal/internal/modules/category"
)

// CategoryHandler handles HTTP requests for categories.
type CategoryHandler struct {
	service *category.Service
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(service *category.Service) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// categoryResponse is the JSON response for a category.
type categoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCategoryResponse(c category.Category) categoryResponse {
	return categoryResponse{
		ID:        c.ID.String(),
		Name:      c.Name,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCategoryResponses(categories []category.Category) []categoryResponse {
	result := make([]categoryResponse, len(categories))
	for i, c := range categories {
		result[i] = toCategoryResponse(c)
	}
	return result
}

type createCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

// Create handles POST /categories
func (h *CategoryHandler) Create(c echo.Context) error {
	var req createCategoryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	cat, err := h.service.Create(c.Request().Context(), category.CreateInput{
		Name: req.Name,
	})
	if err != nil {
		if errors.Is(err, category.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, "category already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create category")
	}

	return c.JSON(http.StatusCreated, toCategoryResponse(cat))
}

// GetAll handles GET /categories
func (h *CategoryHandler) GetAll(c echo.Context) error {
	categories, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get categories")
	}

	if categories == nil {
		categories = []category.Category{}
	}

	return c.JSON(http.StatusOK, toCategoryResponses(categories))
}

// GetByID handles GET /categories/:id
func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := category.ParseID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid category id")
	}

	cat, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, category.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get category")
	}

	return c.JSON(http.StatusOK, toCategoryResponse(cat))
}

type updateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

// Update handles PUT /categories/:id
func (h *CategoryHandler) Update(c echo.Context) error {
	id, err := category.ParseID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid category id")
	}

	var req updateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name is required")
	}

	cat, err := h.service.Update(c.Request().Context(), category.UpdateInput{
		ID:   id,
		Name: req.Name,
	})
	if err != nil {
		if errors.Is(err, category.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update category")
	}

	return c.JSON(http.StatusOK, toCategoryResponse(cat))
}

// Delete handles DELETE /categories/:id
func (h *CategoryHandler) Delete(c echo.Context) error {
	id, err := category.ParseID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid category id")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		if errors.Is(err, category.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete category")
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers category routes.
func (h *CategoryHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
