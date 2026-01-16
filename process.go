package jobtick

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Process status constants
const (
	StatusCreated  = "created"
	StatusAssigned = "assigned"
	StatusExecuted = "executed"
	StatusFailed   = "failed"
)

// Process type constants
const (
	Call      = "1"
	Recurring = "2"
)

// Process represents a background process.
type Process struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Type      string          `json:"type"`
	Queue     string          `json:"queue"`
	File      string          `json:"file"`
	Function  string          `json:"function"`
	Data      json.RawMessage `json:"data"`
	StartTime time.Time       `json:"start_time"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Create inserts a new process into the database.
func (jt *JobTick) Create(name, procType, queue, file, function string, data json.RawMessage, startTime time.Time) error {
	query := `
		INSERT INTO processes (name, type, queue, file, function, data, start_time, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`
	_, err := jt.db.Exec(query, name, procType, queue, file, function, data, startTime, StatusCreated)
	return err
}

// Update updates specific fields of a process.
func (jt *JobTick) Update(id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}

	query := "UPDATE processes SET "
	args := []interface{}{}
	i := 1

	for field, value := range fields {
		query += fmt.Sprintf("%s = $%d, ", field, i)
		args = append(args, value)
		i++
	}

	query += fmt.Sprintf("updated_at = NOW() WHERE id = $%d", i)
	args = append(args, id)

	_, err := jt.db.Exec(query, args...)
	return err
}

// UpdateStatus updates the status of a process.
func (jt *JobTick) UpdateStatus(id string, status string) error {
	return jt.Update(id, map[string]interface{}{
		"status": status,
	})
}

// FromID fetches a process by its ID.
func (jt *JobTick) FromID(id string) (*Process, error) {
	query := `
		SELECT id, name, type, queue, file, function, data, start_time, status, created_at, updated_at
		FROM processes
		WHERE id = $1
	`
	proc := &Process{}
	err := jt.db.QueryRow(query, id).Scan(
		&proc.ID, &proc.Name, &proc.Type, &proc.Queue, &proc.File, &proc.Function,
		&proc.Data, &proc.StartTime, &proc.Status, &proc.CreatedAt, &proc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return proc, nil
}

// GetNextForExecution fetches the next pending process and marks it as assigned.
func (jt *JobTick) GetNextForExecution(queue string) (*Process, error) {
	query := `
		SELECT id, name, type, queue, file, function, data, start_time, status, created_at, updated_at
		FROM processes
		WHERE queue = $1 AND status = $2 AND start_time <= NOW()
		ORDER BY start_time ASC
		LIMIT 1
		FOR UPDATE SKIP LOCKED
	`
	proc := &Process{}
	err := jt.db.QueryRow(query, queue, StatusCreated).Scan(
		&proc.ID, &proc.Name, &proc.Type, &proc.Queue, &proc.File, &proc.Function,
		&proc.Data, &proc.StartTime, &proc.Status, &proc.CreatedAt, &proc.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	err = jt.UpdateStatus(proc.ID, StatusAssigned)
	if err != nil {
		return nil, err
	}

	proc.Status = StatusAssigned
	return proc, nil
}
