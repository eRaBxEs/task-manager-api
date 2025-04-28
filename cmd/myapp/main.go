package main

import (
	"fmt"
	"log"
	"net/http"

	"task-api/internal/config"
	"task-api/internal/database"
	"task-api/internal/routes"
)

func main() {

	//Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config %v", err)
	}

	// Initialize database connection
	db, err := database.InitDB(cfg.DB) //now uses config
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database!")

	//Set up routes
	r := routes.SetupRoutes(db) // Pass db to routes
	//Start server
	log.Fatal(http.ListenAndServe(":8080", r)) // Start the server using http.ListenAndServe
}
