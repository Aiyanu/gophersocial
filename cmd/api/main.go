package main

import (
	"log"
	"time"

	"github.com/Aiyanu/gophersocial/internal/auth"
	"github.com/Aiyanu/gophersocial/internal/db"
	"github.com/Aiyanu/gophersocial/internal/env"
	"github.com/Aiyanu/gophersocial/internal/mailer"
	"github.com/Aiyanu/gophersocial/internal/store"
	"github.com/Aiyanu/gophersocial/internal/store/cache"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://aiyanu:incorrect@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			fromEmail: env.GetString("FROM_EMAIL", "aiyanu1.000@outlook.com"),
			sendGrid: sendgridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			exp: time.Hour * 24 * 3,
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOEKN_SECRET", "example"),
				exp:    time.Hour * 24 * 3,
				iss:    "gophersocial",
			},
		},
	}

	//Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Establish database connection with better error handling
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	defer db.Close()
	logger.Info("database connection pool established")
	var rdb *redis.Client

	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
		logger.Info("Redis cache connection Established")
	}

	// Initialize the store and application
	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)
	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenicator := auth.NewJWTAuthenicator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)
	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenicator,
	}

	// Start the server
	mux := app.mount()
	log.Printf("Starting server on %s in %s mode\n", cfg.addr, cfg.env)
	log.Fatal(app.run(mux))
}
