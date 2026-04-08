package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"userprofile-api/models"
)

// TestDB wraps a GORM DB for testing with proper cleanup.
type TestDB struct {
	DB *gorm.DB
}

// NewTestDB creates an in-memory SQLite database for testing.
func NewTestDB() (*TestDB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&models.UserProfile{}); err != nil {
		return nil, err
	}

	return &TestDB{DB: db}, nil
}

// Close closes the underlying database connection.
func (tdb *TestDB) Close() {
	if sqlDB, err := tdb.DB.DB(); err == nil {
		sqlDB.Close()
	}
}
