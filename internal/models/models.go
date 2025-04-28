package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Task represents a task in the database.
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"` // Must match frontend field name
}

// Custom unmarshal to handle both RFC3339 and simple date formats
func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		DueDate string `json:"due_date"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.DueDate == "" {
		return errors.New("due_date is required")
	}

	// Try parsing as RFC3339 (full ISO format)
	parsedTime, err := time.Parse(time.RFC3339, aux.DueDate)
	if err != nil {
		// Try parsing as simple date if RFC3339 fails
		parsedTime, err = time.Parse("2006-01-02T15:04", aux.DueDate)
		if err != nil {
			return fmt.Errorf("invalid due_date format: %v", err)
		}
	}

	t.DueDate = parsedTime
	return nil
}
