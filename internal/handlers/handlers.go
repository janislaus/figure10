package handlers

import (
	"database/sql"

	"github.com/janislaus/figure10/internal/llm"
)

// Handler holds dependencies for the HTTP handlers
type Handler struct {
	DB        *sql.DB
	Generator *llm.TextGenerator
}

// NewHandler creates a new Handler with the given dependencies
func NewHandler(db *sql.DB, generator *llm.TextGenerator) *Handler {
	return &Handler{
		DB:        db,
		Generator: generator,
	}
}
