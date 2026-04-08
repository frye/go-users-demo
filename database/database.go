package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"userprofile-api/models"
)

// InitDB opens a SQLite database and runs migrations.
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.UserProfile{}); err != nil {
		return nil, err
	}

	return db, nil
}

// SeedDB populates the database with default users if the table is empty.
func SeedDB(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.UserProfile{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		users := []models.UserProfile{
			{FullName: "John Doe", Emoji: "😀"},
			{FullName: "Jane Smith", Emoji: "🚀"},
			{FullName: "Robert Johnson", Emoji: "🎸"},
		}
		if err := db.Create(&users).Error; err != nil {
			return err
		}
	}

	return nil
}
