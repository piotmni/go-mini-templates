package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/piotmni/go-mini-templates/minimal/internal/modules/category"
	"github.com/piotmni/go-mini-templates/minimal/internal/modules/note"
)

// NoteHandler handles HTTP requests for notes.
type NoteHandler struct {
	service *note.Service
}

// NewNoteHandler creates a new NoteHandler.
func NewNoteHandler(service *note.Service) *NoteHandler {
	return &NoteHandler{service: service}
}

// noteResponse is the JSON response for a note.
type noteResponse struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func toNoteResponse(n note.Note) noteResponse {
	return noteResponse{
		ID:         n.ID.String(),
		CategoryID: n.CategoryID.String(),
		Title:      n.Title,
		Content:    n.Content,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}

func toNoteResponses(notes []note.Note) []noteResponse {
	result := make([]noteResponse, len(notes))
	for i, n := range notes {
		result[i] = toNoteResponse(n)
	}
	return result
}

type createNoteRequest struct {
	CategoryID string `json:"category_id" validate:"required"`
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content"`
}

// Create handles POST /notes
func (h *NoteHandler) Create(c echo.Context) error {
	var req createNoteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "title is required")
	}

	categoryID, err := category.ParseID(req.CategoryID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid category_id")
	}

	n, err := h.service.Create(c.Request().Context(), note.CreateInput{
		CategoryID: categoryID,
		Title:      req.Title,
		Content:    req.Content,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create note")
	}

	return c.JSON(http.StatusCreated, toNoteResponse(n))
}

// GetAll handles GET /notes
func (h *NoteHandler) GetAll(c echo.Context) error {
	// Check if category filter is provided
	categoryIDStr := c.QueryParam("category_id")
	if categoryIDStr != "" {
		categoryID, err := category.ParseID(categoryIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid category_id")
		}

		notes, err := h.service.GetByCategory(c.Request().Context(), categoryID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get notes")
		}

		if notes == nil {
			notes = []note.Note{}
		}

		return c.JSON(http.StatusOK, toNoteResponses(notes))
	}

	notes, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get notes")
	}

	if notes == nil {
		notes = []note.Note{}
	}

	return c.JSON(http.StatusOK, toNoteResponses(notes))
}

// GetByID handles GET /notes/:id
func (h *NoteHandler) GetByID(c echo.Context) error {
	id, err := note.ParseID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid note id")
	}

	n, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, note.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "note not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get note")
	}

	return c.JSON(http.StatusOK, toNoteResponse(n))
}

type updateNoteRequest struct {
	CategoryID string `json:"category_id" validate:"required"`
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content"`
}

// Update handles PUT /notes/:id
func (h *NoteHandler) Update(c echo.Context) error {
	id, err := note.ParseID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid note id")
	}

	var req updateNoteRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "title is required")
	}

	categoryID, err := category.ParseID(req.CategoryID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid category_id")
	}

	n, err := h.service.Update(c.Request().Context(), note.UpdateInput{
		ID:         id,
		CategoryID: categoryID,
		Title:      req.Title,
		Content:    req.Content,
	})
	if err != nil {
		if errors.Is(err, note.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "note not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update note")
	}

	return c.JSON(http.StatusOK, toNoteResponse(n))
}

// Delete handles DELETE /notes/:id
func (h *NoteHandler) Delete(c echo.Context) error {
	id, err := note.ParseID(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid note id")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		if errors.Is(err, note.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "note not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete note")
	}

	return c.NoContent(http.StatusNoContent)
}

// RegisterRoutes registers note routes.
func (h *NoteHandler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.Create)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
