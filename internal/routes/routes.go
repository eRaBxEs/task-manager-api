package routes

import (
	"database/sql"
	"net/http"
	"task-api/internal/handlers"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// SetupRoutes sets up the application routes
func SetupRoutes(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	taskHandler := handlers.NewTaskHandler(db)

	// Define your routes
	r.HandleFunc("/tasks", taskHandler.CreateTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks", taskHandler.GetTasks).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods(http.MethodGet)
	r.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods(http.MethodPut)
	r.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods(http.MethodDelete)

	// Apply CORS middleware to all routes
	headersOk := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := gorillaHandlers.AllowedOrigins([]string{"*"}) // or specify your frontend domains instead of "*"
	methodsOk := gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	return gorillaHandlers.CORS(headersOk, originsOk, methodsOk)(r)
}
