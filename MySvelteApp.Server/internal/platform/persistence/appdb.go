package persistence

import (
	"fmt"

	"gorm.io/gorm"

	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

// AppDB wraps gorm.DB to keep persistence wiring centralised.
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

// AutoMigrate applies the schema required for the modules currently in use.
func (a *AppDB) AutoMigrate() error {
	return a.DB.AutoMigrate(&authdomain.User{})
}
