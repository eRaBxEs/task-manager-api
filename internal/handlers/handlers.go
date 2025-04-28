package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"task-api/internal/models"

	"github.com/gorilla/mux"
)

// TaskHandler holds the database connection.
type TaskHandler struct {
	DB *sql.DB
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// CreateTask handles the creation of a new task.
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// First read the raw body for debugging
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	log.Printf("Raw request body: %s", string(bodyBytes))

	// Reset the body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Decoding error: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// The custom unmarshal will have already validated the date
	// Now just verify other required fields
	if task.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if task.Status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	// Database insertion
	var id int
	err = h.DB.QueryRow(
		`INSERT INTO tasks (title, description, status, due_date) 
         VALUES ($1, $2, $3, $4) RETURNING id`,
		task.Title, task.Description, task.Status, task.DueDate,
	).Scan(&id)

	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	var createdTask models.Task
	err = h.DB.QueryRow(
		`SELECT id, title, description, status, due_date 
         FROM tasks WHERE id = $1`, id,
	).Scan(&createdTask.ID, &createdTask.Title, &createdTask.Description,
		&createdTask.Status, &createdTask.DueDate)

	if err != nil {
		log.Printf("Error fetching created task: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask) // Send back the complete task
}

// GetTasks retrieves all tasks.
// GetTasks retrieves all tasks.
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, title, description, status, due_date FROM tasks")
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDate)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}

	// Always return an array, even if empty
	if tasks == nil {
		tasks = []models.Task{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTask retrieves a single task by ID.
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Use the database connection from the handler.
	row := h.DB.QueryRow("SELECT id, title, description, status, due_date FROM tasks WHERE id = $1", id)
	var task models.Task
	err = row.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.DueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		log.Printf("Error scanning row: %v", err)
		http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTask updates an existing task.
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//update the id
	task.ID = id

	// Use the database connection from the handler.
	result, err := h.DB.Exec("UPDATE tasks SET title = $1, description = $2, status = $3, due_date = $4 WHERE id = $5", task.Title, task.Description, task.Status, task.DueDate, task.ID)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected row count: %v", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Updated %d row(s).\n", rowsAffected)
	w.WriteHeader(http.StatusNoContent) //204
}

// DeleteTask deletes a task.
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	result, err := h.DB.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected row count: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	fmt.Printf("Deleted %d row(s).\n", rowsAffected)
	w.WriteHeader(http.StatusNoContent) // 204
}
