package main

import (
	"log"

	"github.com/gophersocial/internal/db"
	"github.com/gophersocial/internal/env"
	"github.com/gophersocial/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {
	// Load environment variables first so we can use them for configuration
	err := godotenv.Load()
	if err != nil {
		// Use a more descriptive error message to help with debugging
		log.Printf("Warning: Error loading .env file: %v. Using default values.\n", err)
	}

	// Create the configuration with the correct Windows host IP
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			// Replace YOUR_WINDOWS_IP with the IP from the nameserver command
			addr:         env.GetString("DB_ADDR", "postgres://aiyanu:incorrect@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Establish database connection with better error handling
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	// Initialize the store and application
	store := store.NewStorage(db)
	app := &application{
		config: cfg,
		store:  store,
	}

	// Start the server
	mux := app.mount()
	log.Printf("Starting server on %s in %s mode\n", cfg.addr, cfg.env)
	log.Fatal(app.run(mux))
}
