package http

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/piotmni/go-mini-templates/ent-atlas/internal/http/handler"
	"github.com/piotmni/go-mini-templates/ent-atlas/internal/http/middleware"
	"github.com/piotmni/go-mini-templates/ent-atlas/internal/service"
)

// RouterDeps holds all dependencies needed by the router.
type RouterDeps struct {
	PasteService *service.PasteService
}

// NewRouter creates a new HTTP router with all routes configured.
func NewRouter(deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes - Paste
	pasteHandler := handler.NewPasteHandler(deps.PasteService)
	mux.HandleFunc("POST /pastes", pasteHandler.Create)
	mux.HandleFunc("GET /pastes", pasteHandler.List)
	mux.HandleFunc("GET /pastes/{slug}", pasteHandler.Get)
	mux.HandleFunc("DELETE /pastes/{slug}", pasteHandler.Delete)

	// Static files
	staticDir := getStaticDir()
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Serve index.html for root and paste view routes
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	mux.HandleFunc("GET /p/{slug}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.Logger(handler)

	return handler
}

// getStaticDir returns the path to the static directory.
func getStaticDir() string {
	// Check if running from project root
	if _, err := os.Stat("static"); err == nil {
		return "static"
	}
	// Check relative to binary location
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), "static")
}
