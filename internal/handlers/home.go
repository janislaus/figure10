package handlers

import (
	"context"
	"net/http"

	"github.com/janislaus/figure10/internal/db"
	"github.com/janislaus/figure10/web/templates"
)

// HandleHome renders the home page
func (h *Handler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Render the home template
	templates.Base(templates.Home()).Render(context.Background(), w)
}

// HandleHistory renders the history page
func (h *Handler) HandleHistory(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get recent sessions
	sessions, err := db.GetRecentSessions(h.DB, 10)
	if err != nil {
		http.Error(w, "Failed to load history", http.StatusInternalServerError)
		return
	}

	// Get common errors
	errors, err := db.GetCommonErrors(h.DB, 10)
	if err != nil {
		http.Error(w, "Failed to load common errors", http.StatusInternalServerError)
		return
	}

	// Render the history template
	templates.Base(templates.History(sessions, errors)).Render(ctx, w)
}
