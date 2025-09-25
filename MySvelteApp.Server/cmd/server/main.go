package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mysvelteapp/server_new/internal/docs"
	authapi "mysvelteapp/server_new/internal/modules/auth/api"
	authapp "mysvelteapp/server_new/internal/modules/auth/app"
	authpersistence "mysvelteapp/server_new/internal/modules/auth/infra/persistence"
	authsecurity "mysvelteapp/server_new/internal/modules/auth/infra/security"
	authtoken "mysvelteapp/server_new/internal/modules/auth/infra/token"
	pokemonapi "mysvelteapp/server_new/internal/modules/pokemon/api"
	pokemonapp "mysvelteapp/server_new/internal/modules/pokemon/app"
	pokemoninfra "mysvelteapp/server_new/internal/modules/pokemon/infra/pokeapi"
	"mysvelteapp/server_new/internal/platform/config"
	"mysvelteapp/server_new/internal/platform/httpserver"
	"mysvelteapp/server_new/internal/platform/logging"
	"mysvelteapp/server_new/internal/platform/persistence"
	"mysvelteapp/server_new/internal/platform/tracing"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logging.NewDefaultLogger()

	// Initialize OpenTelemetry tracing
	tracingProvider, err := tracing.New(cfg.ServiceName, cfg.ServiceVersion, logger)
	if err != nil {
		log.Fatalf("failed to initialize tracing: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracingProvider.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown tracing provider: %v", err)
		}
	}()

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "MySvelteApp Server API"
	docs.SwaggerInfo.Description = "This is the Go implementation of the MySvelteApp backend."

	engine := httpserver.New(logger, cfg.ServiceName)

	appDB, err := persistence.NewAppDB(sqlite.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}
	if err := appDB.AutoMigrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	passwordHasher := authsecurity.NewHMACPasswordHasher()

	tokenGenerator, err := authtoken.NewJWTTokenGenerator(authtoken.JWTOptions{
		Key:                      cfg.JWTKey,
		Issuer:                   cfg.JWTIssuer,
		Audience:                 cfg.JWTAudience,
		AccessTokenLifetimeHours: cfg.JWTAccessLifetimeHours,
	})
	if err != nil {
		log.Fatalf("failed to initialise JWT generator: %v", err)
	}

	userRepository := authpersistence.NewGormUserRepository(appDB.DB)
	authService := authapp.NewService(userRepository, passwordHasher, tokenGenerator)
	authHandlers := authapi.NewHandlers(authService)
	authapi.RegisterRoutes(engine, authHandlers)

	pokemonAdapter := pokemoninfra.NewAdapter(http.DefaultClient)
	pokemonService := pokemonapp.NewService(pokemonAdapter)
	pokemonHandlers := pokemonapi.NewHandlers(pokemonService)
	pokemonapi.RegisterRoutes(engine, pokemonHandlers)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: engine,
	}

	go func() {
		log.Printf("Server listening on http://localhost:%s", cfg.Port)
		log.Printf("OpenTelemetry tracing enabled (development mode: stdout exporter)")
		log.Printf("Service: %s v%s", cfg.ServiceName, cfg.ServiceVersion)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
