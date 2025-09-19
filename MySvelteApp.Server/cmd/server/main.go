package main

import (
	"log"
	"net/http"

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
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := logging.NewDefaultLogger()

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "MySvelteApp Server API"
	docs.SwaggerInfo.Description = "This is the Go implementation of the MySvelteApp backend."

	engine := httpserver.New(logger)

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

	log.Printf("Server listening on http://localhost:%s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
