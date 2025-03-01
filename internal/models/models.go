package models

import "time"

// Text represents a typing exercise text
type Text struct {
	ID        int64
	Content   string
	Prompt    string
	CreatedAt time.Time
}

// Session represents a typing session
type Session struct {
	ID          int64
	TextID      int64
	WPM         float64
	Accuracy    float64
	Errors      int
	CompletedAt time.Time
}

// SessionWithText extends Session with the text prompt
type SessionWithText struct {
	Session
	Prompt string
}

// TypingError represents a specific typing error
type TypingError struct {
	ID           int64
	SessionID    int64
	ExpectedChar string
	TypedChar    string
	Position     int
}

// CommonError represents a common typing error
type CommonError struct {
	ExpectedChar string
	TypedChar    string
	Count        int
}

// TypingResult represents the result of a typing session
type TypingResult struct {
	TextID       int64         `json:"text_id"`
	WPM          float64       `json:"wpm"`
	Accuracy     float64       `json:"accuracy"`
	Errors       int           `json:"errors"`
	ErrorDetails []TypingError `json:"error_details"`
	ErrorWords   []string      `json:"error_words"`
}

// TypingCheck represents a real-time typing check result
type TypingCheck struct {
	Correct    bool
	CurrentWPM float64
	CurrentAcc float64
	CurrentPos int
	TotalChars int
	ErrorCount int
}
