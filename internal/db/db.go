package db

import (
	"database/sql"
	"time"

	"github.com/janislaus/figure10/internal/models"
)

// InitDB initializes the database schema
func InitDB(db *sql.DB) error {
	// Create texts table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS texts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			prompt TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create sessions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text_id INTEGER,
			wpm REAL,
			accuracy REAL,
			errors INTEGER,
			completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (text_id) REFERENCES texts(id)
		)
	`)
	if err != nil {
		return err
	}

	// Create errors table to track specific errors
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS typing_errors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER,
			expected_char TEXT,
			typed_char TEXT,
			position INTEGER,
			FOREIGN KEY (session_id) REFERENCES sessions(id)
		)
	`)

	return err
}

// SaveText saves a new text to the database
func SaveText(db *sql.DB, content, prompt string) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO texts (content, prompt) VALUES (?, ?)",
		content, prompt,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetTextByID retrieves a text by its ID
func GetTextByID(db *sql.DB, id int64) (models.Text, error) {
	var text models.Text
	var createdAtStr string

	err := db.QueryRow(
		"SELECT id, content, prompt, created_at FROM texts WHERE id = ?",
		id,
	).Scan(&text.ID, &text.Content, &text.Prompt, &createdAtStr)

	if err != nil {
		return models.Text{}, err
	}

	text.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	return text, nil
}

// SaveSession saves a new typing session to the database
func SaveSession(db *sql.DB, textID int64, wpm, accuracy float64, errors int) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO sessions (text_id, wpm, accuracy, errors) VALUES (?, ?, ?, ?)",
		textID, wpm, accuracy, errors,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// SaveTypingError saves a typing error to the database
func SaveTypingError(db *sql.DB, sessionID int64, expected, typed string, position int) error {
	_, err := db.Exec(
		"INSERT INTO typing_errors (session_id, expected_char, typed_char, position) VALUES (?, ?, ?, ?)",
		sessionID, expected, typed, position,
	)
	return err
}

// GetRecentSessions retrieves recent typing sessions
func GetRecentSessions(db *sql.DB, limit int) ([]models.SessionWithText, error) {
	rows, err := db.Query(`
		SELECT s.id, s.text_id, s.wpm, s.accuracy, s.errors, s.completed_at, t.prompt 
		FROM sessions s
		JOIN texts t ON s.text_id = t.id
		ORDER BY s.completed_at DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.SessionWithText
	for rows.Next() {
		var session models.SessionWithText
		var completedAtStr string

		err := rows.Scan(
			&session.ID,
			&session.TextID,
			&session.WPM,
			&session.Accuracy,
			&session.Errors,
			&completedAtStr,
			&session.Prompt,
		)

		if err != nil {
			return nil, err
		}

		session.CompletedAt, _ = time.Parse("2006-01-02 15:04:05", completedAtStr)
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetCommonErrors retrieves the most common typing errors
func GetCommonErrors(db *sql.DB, limit int) ([]models.CommonError, error) {
	rows, err := db.Query(`
		SELECT expected_char, typed_char, COUNT(*) as count
		FROM typing_errors
		GROUP BY expected_char, typed_char
		ORDER BY count DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var errors []models.CommonError
	for rows.Next() {
		var e models.CommonError

		err := rows.Scan(&e.ExpectedChar, &e.TypedChar, &e.Count)
		if err != nil {
			return nil, err
		}

		errors = append(errors, e)
	}

	return errors, nil
}
