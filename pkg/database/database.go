package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseService handles all database operations
type DatabaseService struct {
	db *gorm.DB
}

// NewDatabaseService creates a new database service
// Models should be passed in to enable auto-migration
func NewDatabaseService(databaseURL string, models ...interface{}) (*DatabaseService, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema if models are provided
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
	}

	return &DatabaseService{db: db}, nil
}

// Create creates a new record in the database
func (s *DatabaseService) Create(model interface{}) error {
	if err := s.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create record: %w", err)
	}
	return nil
}

// GetByID retrieves a record by its ID
func (s *DatabaseService) GetByID(model interface{}, id string) error {
	if err := s.db.First(model, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to get record: %w", err)
	}
	return nil
}

// List retrieves all records with pagination
func (s *DatabaseService) List(model interface{}, result interface{}, limit, offset int) (int64, error) {
	var total int64

	// Get total count
	if err := s.db.Model(model).Count(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to count records: %w", err)
	}

	// Get paginated results
	if err := s.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(result).Error; err != nil {
		return 0, fmt.Errorf("failed to list records: %w", err)
	}

	return total, nil
}

// Update updates an existing record
func (s *DatabaseService) Update(model interface{}, id string, updates map[string]interface{}) error {
	// Find the record
	if err := s.db.First(model, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to find record: %w", err)
	}

	// Update the record
	if err := s.db.Model(model).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	return nil
}

// Delete deletes a record by its ID
func (s *DatabaseService) Delete(model interface{}, id string) error {
	result := s.db.Delete(model, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

// GetDB returns the underlying GORM database instance for custom queries
func (s *DatabaseService) GetDB() *gorm.DB {
	return s.db
}

// Close closes the database connection
func (s *DatabaseService) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
