package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/janislaus/figure10/internal/db"
	"github.com/janislaus/figure10/internal/models"
	"github.com/janislaus/figure10/web/templates"
)

// HandleGenerateText generates a new typing text
func (h *Handler) HandleGenerateText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the prompt from the form
	prompt := r.FormValue("prompt")
	if prompt == "" {
		prompt = "Give me a general typing practice text"
	}

	// Generate text using the LLM
	content, err := h.Generator.GenerateText(prompt)
	if err != nil {
		http.Error(w, "Failed to generate text", http.StatusInternalServerError)
		return
	}

	// Save the text to the database
	textID, err := db.SaveText(h.DB, content, prompt)
	if err != nil {
		http.Error(w, "Failed to save text", http.StatusInternalServerError)
		return
	}

	// Create a Text model
	text := models.Text{
		ID:        textID,
		Content:   content,
		Prompt:    prompt,
		CreatedAt: time.Now(),
	}

	// Render the typing exercise template
	templates.TypingExercise(text).Render(context.Background(), w)
}

// HandleStartSession starts a new typing session
func (h *Handler) HandleStartSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the text ID from the form
	textIDStr := r.FormValue("text_id")
	textID, err := strconv.ParseInt(textIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid text ID", http.StatusBadRequest)
		return
	}

	// Get the text from the database
	text, err := db.GetTextByID(h.DB, textID)
	if err != nil {
		http.Error(w, "Failed to get text", http.StatusInternalServerError)
		return
	}

	// Return the text content as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"text_id": text.ID,
		"content": text.Content,
		"prompt":  text.Prompt,
	})
}

// HandleCheckTyping checks the current typing progress
func (h *Handler) HandleCheckTyping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the text ID, current input, and start time from the form
	textIDStr := r.FormValue("text_id")
	textID, err := strconv.ParseInt(textIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid text ID", http.StatusBadRequest)
		return
	}

	currentInput := r.FormValue("current_input")
	startTimeStr := r.FormValue("start_time")
	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid start time", http.StatusBadRequest)
		return
	}

	// Get the text from the database
	text, err := db.GetTextByID(h.DB, textID)
	if err != nil {
		http.Error(w, "Failed to get text", http.StatusInternalServerError)
		return
	}

	// Calculate metrics
	elapsedTime := float64(time.Now().UnixMilli()-startTime) / 1000.0 // in seconds
	totalChars := len(text.Content)
	currentPos := len(currentInput)

	// Count errors
	errorCount := 0
	for i := 0; i < currentPos && i < totalChars; i++ {
		if i >= len(currentInput) || text.Content[i] != currentInput[i] {
			errorCount++
		}
	}

	// Calculate WPM (assuming 5 chars per word)
	var wpm float64
	if elapsedTime > 0 {
		wpm = float64(currentPos) / 5.0 / (elapsedTime / 60.0)
	}

	// Calculate accuracy
	var accuracy float64
	if currentPos > 0 {
		accuracy = 100.0 * float64(currentPos-errorCount) / float64(currentPos)
	}

	// Check if the current character is correct
	correct := true
	if currentPos < totalChars && currentPos > 0 {
		if currentInput[currentPos-1] != text.Content[currentPos-1] {
			correct = false
		}
	}

	// Create the response
	check := models.TypingCheck{
		Correct:    correct,
		CurrentWPM: wpm,
		CurrentAcc: accuracy,
		CurrentPos: currentPos,
		TotalChars: totalChars,
		ErrorCount: errorCount,
	}

	// Return the check result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(check)
}

// HandleSubmitResult submits the final result of a typing session
func (h *Handler) HandleSubmitResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var result models.TypingResult
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Save the session to the database
	sessionID, err := db.SaveSession(h.DB, result.TextID, result.WPM, result.Accuracy, result.Errors)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Save the error details
	for _, e := range result.ErrorDetails {
		err := db.SaveTypingError(h.DB, sessionID, e.ExpectedChar, e.TypedChar, e.Position)
		if err != nil {
			// Log the error but continue
			fmt.Printf("Failed to save typing error: %v\n", err)
		}
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"session_id": sessionID,
	})
}

// HandleGeneratePractice generates a practice text with words that had errors
func (h *Handler) HandleGeneratePractice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var request struct {
		Words []string `json:"words"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(request.Words) == 0 {
		http.Error(w, "No words provided", http.StatusBadRequest)
		return
	}

	// Create a more specific prompt that ensures each word appears multiple times
	prompt := fmt.Sprintf(
		"Create a typing practice paragraph that includes EACH of these words AT LEAST 10 TIMES: %s. "+
			"Make sure each word appears multiple times throughout the text. "+
			"The text should be coherent but focus on repeating these words frequently for practice.",
		strings.Join(request.Words, ", "))

	fmt.Printf("Generating practice with prompt: %s\n", prompt)

	// Generate text using the LLM
	content, err := h.Generator.GenerateText(prompt)
	if err != nil {
		http.Error(w, "Failed to generate practice text", http.StatusInternalServerError)
		return
	}

	// Verify that each word appears multiple times
	for _, word := range request.Words {
		count := strings.Count(strings.ToLower(content), strings.ToLower(word))
		fmt.Printf("Word '%s' appears %d times in generated text\n", word, count)
	}

	// Save the text to the database
	textID, err := db.SaveText(h.DB, content, "Practice: "+strings.Join(request.Words, ", "))
	if err != nil {
		http.Error(w, "Failed to save text", http.StatusInternalServerError)
		return
	}

	// Create a Text model
	text := models.Text{
		ID:        textID,
		Content:   content,
		Prompt:    prompt,
		CreatedAt: time.Now(),
	}

	// Render the typing exercise template
	templates.TypingExercise(text).Render(context.Background(), w)
}
