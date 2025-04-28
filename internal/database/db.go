package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"task-api/internal/config"
)

// InitDB initializes the database connection
func InitDB(cfg config.DBConfig) (*sql.DB, error) {
	//connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
	//  os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Name, cfg.User, cfg.Password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close() //important to close db here
		return nil, fmt.Errorf("error pinging the database: %w", err)
	}
	return db, nil
}
