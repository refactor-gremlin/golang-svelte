package persistence

import (
	"fmt"

	"gorm.io/gorm"

	"mysvelteapp/server/internal/domain/entities"
)

// AppDB wraps gorm.DB to keep the infrastructure layer cohesive.
type AppDB struct {
	DB *gorm.DB
}

// NewAppDB constructs an AppDB given a prepared gorm dialector and configuration.
func NewAppDB(dialector gorm.Dialector, config *gorm.Config) (*AppDB, error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	return &AppDB{DB: db}, nil
}

// AutoMigrate applies the schema required for the auth module.
func (a *AppDB) AutoMigrate() error {
	return a.DB.AutoMigrate(&entities.User{})
}
