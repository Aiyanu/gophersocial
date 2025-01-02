package main

import (
	"log"

	"github.com/Aiyanu/gophersocial/internal/db"
	"github.com/Aiyanu/gophersocial/internal/env"
	"github.com/Aiyanu/gophersocial/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description

func main() {
	// Load environment variables first so we can use them for configuration
	err := godotenv.Load()
	if err != nil {
		// Use a more descriptive error message to help with debugging
		log.Printf("Warning: Error loading .env file: %v. Using default values.\n", err)
	}

	// Create the configuration with the correct Windows host IP
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "http://localhost:8080"),
		db: dbConfig{
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
