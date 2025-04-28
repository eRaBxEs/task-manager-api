package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"task-api/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// testDB wraps both the *sql.DB and sqlmock for testing
type testDB struct {
	DB   *sql.DB
	mock sqlmock.Sqlmock
}

func newTestDB(t *testing.T) *testDB {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	return &testDB{DB: db, mock: mock}
}

func TestCreateTask(t *testing.T) {
	// Setup mock DB
	tDB := newTestDB(t)
	defer tDB.DB.Close()

	// Create handler with the mocked *sql.DB
	handler := NewTaskHandler(tDB.DB)

	// Test cases
	tests := []struct {
		name           string
		payload        interface{}
		setupMock      func()
		expectedStatus int
		expectedTask   models.Task
	}{
		{
			name: "successful task creation",
			payload: models.Task{
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				DueDate:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			setupMock: func() {
				// Expect the INSERT
				tDB.mock.ExpectExec(`INSERT INTO tasks \(title, description, status, due_date\) VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
					WithArgs("Test Task", "Test Description", "pending", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Expect the SELECT to return the created task
				tDB.mock.ExpectQuery(`SELECT id, title, description, status, due_date FROM tasks WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "due_date"}).
						AddRow(1, "Test Task", "Test Description", "pending", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)))
			},
			expectedStatus: http.StatusCreated,
			expectedTask: models.Task{
				ID:          1,
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "pending",
				DueDate:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "invalid payload",
			payload: map[string]interface{}{
				"title": 123, // Invalid type for title
			},
			setupMock:      func() {}, // No DB expectations
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock expectations
			tt.setupMock()

			// Create test request
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			// Call the handler
			handler.CreateTask(rr, req)

			// Verify HTTP response
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// For successful creation, verify response body
			if tt.expectedStatus == http.StatusCreated {
				var response struct {
					Message string      `json:"message"`
					Task    models.Task `json:"task"`
				}
				err := json.NewDecoder(rr.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, "Task created successfully", response.Message)
				assert.Equal(t, tt.expectedTask, response.Task)
			}

			// Verify all mock expectations were met
			assert.NoError(t, tDB.mock.ExpectationsWereMet())
		})
	}
}

func TestGetTask(t *testing.T) {
	tDB := newTestDB(t)
	defer tDB.DB.Close()

	handler := NewTaskHandler(tDB.DB)

	// Set up test data in mock
	tDB.mock.ExpectQuery("SELECT id, title, description, status, due_date FROM tasks WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "status", "due_date"}).
			AddRow(1, "Test Task", "Test Desc", "pending", time.Now()))

	req := httptest.NewRequest("GET", "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rr := httptest.NewRecorder()

	handler.GetTask(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, tDB.mock.ExpectationsWereMet())
}
