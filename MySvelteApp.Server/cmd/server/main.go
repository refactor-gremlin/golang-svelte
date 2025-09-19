package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	appauth "mysvelteapp/server/internal/application/auth"
	_ "mysvelteapp/server/internal/docs"
	"mysvelteapp/server/internal/infrastructure/authentication"
	"mysvelteapp/server/internal/infrastructure/external"
	"mysvelteapp/server/internal/infrastructure/logging"
	"mysvelteapp/server/internal/infrastructure/middleware"
	"mysvelteapp/server/internal/infrastructure/persistence"
	"mysvelteapp/server/internal/infrastructure/security"
	"mysvelteapp/server/internal/presentation/controllers"
)

const (
	defaultPort             = "8080"
	defaultDatabaseDSN      = "file:mysvelteapp.db?cache=shared&_fk=1"
	defaultJWTKey           = "base64:YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE=" // base64-encoded 32 byte key
	defaultJWTIssuer        = "mysvelteapp"
	defaultJWTAudience      = "mysvelteapp"
	defaultJWTLifetimeHours = 24
)

// @title MySvelteApp Server API
// @version 1.0
// @description This is the Go implementation of the MySvelteApp backend.
// @BasePath /
func main() {
	port := getEnv("SERVER_PORT", defaultPort)
	dsn := getEnv("DATABASE_DSN", defaultDatabaseDSN)

	jwtKey := getEnv("JWT_KEY", defaultJWTKey)
	jwtIssuer := getEnv("JWT_ISSUER", defaultJWTIssuer)
	jwtAudience := getEnv("JWT_AUDIENCE", defaultJWTAudience)

	jwtLifetime := defaultJWTLifetimeHours
	if hoursStr := os.Getenv("JWT_ACCESS_TOKEN_LIFETIME_HOURS"); hoursStr != "" {
		if hours, err := strconv.Atoi(hoursStr); err == nil {
			jwtLifetime = hours
		}
	}

	appDB, err := persistence.NewAppDB(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}

	if err := appDB.AutoMigrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	passwordHasher := security.NewHMACPasswordHasher()

	jwtOptions := authentication.JWTOptions{
		Key:                      jwtKey,
		Issuer:                   jwtIssuer,
		Audience:                 jwtAudience,
		AccessTokenLifetimeHours: jwtLifetime,
	}

	tokenGenerator, err := authentication.NewJWTTokenGenerator(jwtOptions)
	if err != nil {
		log.Fatalf("failed to initialise JWT generator: %v", err)
	}

	userRepository := persistence.NewUserRepository(appDB.DB)

	authService := appauth.NewService(userRepository, passwordHasher, tokenGenerator)

	// Initialize Pokemon service and controller
	pokemonService := external.NewPokeApiRandomPokemonService(nil) // nil uses default HTTP client
	pokemonController := controllers.NewRandomPokemonController(pokemonService)

	authController := controllers.NewAuthController(authService)

	// Initialize logger
	appLogger := logging.NewDefaultLogger()

	// Create logging middleware
	loggingMiddleware := middleware.NewLoggingMiddleware(appLogger)

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", authController.Register)
	mux.HandleFunc("/auth/login", authController.Login)
	mux.HandleFunc("/RandomPokemon", pokemonController.Get)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Apply logging middleware to all routes
	loggedMux := loggingMiddleware.Middleware(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggedMux,
	}

	log.Printf("Server listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
