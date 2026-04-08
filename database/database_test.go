package database

import (
	"testing"

	"userprofile-api/models"
)

func setupTestDB(t *testing.T) *TestDB {
	t.Helper()
	tdb, err := NewTestDB()
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	return tdb
}

func TestInitDB(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.Close()

	if tdb.DB == nil {
		t.Fatal("expected non-nil database")
	}

	// Verify the table was created by inserting a record
	user := models.UserProfile{FullName: "Test User", Emoji: "🧪"}
	if err := tdb.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to insert into migrated table: %v", err)
	}

	if user.ID == 0 {
		t.Error("expected auto-generated ID, got 0")
	}
}

func TestSeedDB(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.Close()

	if err := SeedDB(tdb.DB); err != nil {
		t.Fatalf("SeedDB failed: %v", err)
	}

	var users []models.UserProfile
	tdb.DB.Find(&users)

	if len(users) != 3 {
		t.Fatalf("expected 3 seeded users, got %d", len(users))
	}

	if users[0].FullName != "John Doe" {
		t.Errorf("expected first user to be John Doe, got %s", users[0].FullName)
	}
}

func TestSeedDBIdempotent(t *testing.T) {
	tdb := setupTestDB(t)
	defer tdb.Close()

	if err := SeedDB(tdb.DB); err != nil {
		t.Fatalf("first SeedDB failed: %v", err)
	}
	if err := SeedDB(tdb.DB); err != nil {
		t.Fatalf("second SeedDB failed: %v", err)
	}

	var count int64
	tdb.DB.Model(&models.UserProfile{}).Count(&count)

	if count != 3 {
		t.Errorf("expected 3 users after double seed, got %d", count)
	}
}
