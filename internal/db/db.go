// db/db.go
package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	// Add connection status logging
	log.Printf("Attempting to connect to PostgreSQL at: %s", addr)

	db, err := sql.Open("postgres", addr)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		log.Printf("Failed to parse idle time duration: %v", err)
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// Create context with timeout and add logging
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Attempting to ping database...")
	if err = db.PingContext(ctx); err != nil {
		log.Printf("Ping failed: %v", err)
		return nil, err
	}

	log.Printf("Successfully connected to database")
	return db, nil
}
