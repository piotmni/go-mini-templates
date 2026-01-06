package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/piotmni/go-mini-templates/ent-atlas/internal/service"
)

// PasteHandler handles paste HTTP requests.
type PasteHandler struct {
	service *service.PasteService
}

// NewPasteHandler creates a new paste handler.
func NewPasteHandler(s *service.PasteService) *PasteHandler {
	return &PasteHandler{service: s}
}

// CreateRequest represents the create paste request body.
type CreateRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Language  string `json:"language"`
	IsPublic  *bool  `json:"is_public"`
	ExpiresIn *int   `json:"expires_in_minutes"` // minutes
}

// Create handles POST /pastes
func (h *PasteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	input := service.CreatePasteInput{
		Title:    req.Title,
		Content:  req.Content,
		Language: req.Language,
		IsPublic: true,
	}

	if req.IsPublic != nil {
		input.IsPublic = *req.IsPublic
	}

	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		duration := time.Duration(*req.ExpiresIn) * time.Minute
		input.ExpiresIn = &duration
	}

	paste, err := h.service.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create paste")
		return
	}

	writeJSON(w, http.StatusCreated, paste)
}

// Get handles GET /pastes/{slug}
func (h *PasteHandler) Get(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	paste, err := h.service.GetBySlug(r.Context(), slug)
	if err != nil {
		if _, ok := err.(*service.ErrPasteExpired); ok {
			writeError(w, http.StatusGone, "paste has expired")
			return
		}
		writeError(w, http.StatusNotFound, "paste not found")
		return
	}

	writeJSON(w, http.StatusOK, paste)
}

// List handles GET /pastes
func (h *PasteHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	pastes, err := h.service.List(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list pastes")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"pastes": pastes,
		"limit":  limit,
		"offset": offset,
	})
}

// Delete handles DELETE /pastes/{slug}
func (h *PasteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		writeError(w, http.StatusBadRequest, "slug is required")
		return
	}

	if err := h.service.Delete(r.Context(), slug); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete paste")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
