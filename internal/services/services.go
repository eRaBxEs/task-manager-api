package services

import (
	"database/sql"
	"fmt"
	"task-api/internal/models"
)

// TaskService defines the interface for task-related business logic.
type TaskService interface {
	CreateTask(task models.Task) (models.Task, error)
	GetTasks() ([]models.Task, error)
	GetTask(id int) (models.Task, error)
	UpdateTask(task models.Task) error
	DeleteTask(id int) error
}

// taskService is an implementation of TaskService.
type taskService struct {
	db *sql.DB
}

// NewTaskService creates a new TaskService.
func NewTaskService(db *sql.DB) TaskService {
	return &taskService{db: db}
}

// ErrTaskNotFound is returned when a task is not found.
var ErrTaskNotFound = fmt.Errorf("task not found")

// CreateTask creates a new task.
func (s *taskService) CreateTask(task models.Task) (models.Task, error) {
	var id int
	err := s.db.QueryRow(
		"INSERT INTO tasks (title, description, status, due_date) VALUES ($1, $2, $3, $4) RETURNING id",
		task.Title, task.Description, task.Status, task.DueDate,
	).Scan(&id)

	if err != nil {
		return models.Task{}, fmt.Errorf("error inserting task: %w", err)
	}

	task.ID = id
	return task, nil
}

// GetTasks retrieves all tasks.
func (s *taskService) GetTasks() ([]models.Task, error) {
	rows, err := s.db.Query("SELECT id, title, description, status, due_date FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDate)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}
	return tasks, nil
}

// GetTask retrieves a single task by ID.
func (s *taskService) GetTask(id int) (models.Task, error) {
	row := s.db.QueryRow("SELECT id, title, description, status, due_date FROM tasks WHERE id = $1", id)
	var task models.Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Task{}, ErrTaskNotFound
		}
		return models.Task{}, fmt.Errorf("error scanning row: %w", err)
	}
	return task, nil
}

// UpdateTask updates an existing task.
func (s *taskService) UpdateTask(task models.Task) error {
	result, err := s.db.Exec("UPDATE tasks SET title = $1, description = $2, status = $3, due_date = $4 WHERE id = $5", task.Title, task.Description, task.Status, task.DueDate, task.ID)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected row count: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// DeleteTask deletes a task.
func (s *taskService) DeleteTask(id int) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected row count: %w", err)
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
