package models

import (
	"encoding/json"
	"testing"
)

func TestUserProfile(t *testing.T) {
	// Test UserProfile struct creation
	t.Run("CreateUserProfile", func(t *testing.T) {
		user := UserProfile{
			ID:       "1",
			FullName: "John Doe",
			Emoji:    "😀",
		}

		if user.ID != "1" {
			t.Errorf("Expected ID '1', got '%s'", user.ID)
		}
		if user.FullName != "John Doe" {
			t.Errorf("Expected FullName 'John Doe', got '%s'", user.FullName)
		}
		if user.Emoji != "😀" {
			t.Errorf("Expected Emoji '😀', got '%s'", user.Emoji)
		}
	})

	// Test JSON marshaling
	t.Run("JSONMarshaling", func(t *testing.T) {
		user := UserProfile{
			ID:       "2",
			FullName: "Jane Smith",
			Emoji:    "🚀",
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Errorf("Failed to marshal UserProfile to JSON: %v", err)
		}

		expectedJSON := `{"id":"2","fullName":"Jane Smith","emoji":"🚀"}`
		if string(jsonData) != expectedJSON {
			t.Errorf("Expected JSON %s, got %s", expectedJSON, string(jsonData))
		}
	})

	// Test JSON unmarshaling
	t.Run("JSONUnmarshaling", func(t *testing.T) {
		jsonData := `{"id":"3","fullName":"Robert Johnson","emoji":"🎸"}`
		
		var user UserProfile
		err := json.Unmarshal([]byte(jsonData), &user)
		if err != nil {
			t.Errorf("Failed to unmarshal JSON to UserProfile: %v", err)
		}

		if user.ID != "3" {
			t.Errorf("Expected ID '3', got '%s'", user.ID)
		}
		if user.FullName != "Robert Johnson" {
			t.Errorf("Expected FullName 'Robert Johnson', got '%s'", user.FullName)
		}
		if user.Emoji != "🎸" {
			t.Errorf("Expected Emoji '🎸', got '%s'", user.Emoji)
		}
	})

	// Test empty UserProfile
	t.Run("EmptyUserProfile", func(t *testing.T) {
		user := UserProfile{}

		if user.ID != "" {
			t.Errorf("Expected empty ID, got '%s'", user.ID)
		}
		if user.FullName != "" {
			t.Errorf("Expected empty FullName, got '%s'", user.FullName)
		}
		if user.Emoji != "" {
			t.Errorf("Expected empty Emoji, got '%s'", user.Emoji)
		}
	})

	// Test JSON tags
	t.Run("JSONTags", func(t *testing.T) {
		// Test that the JSON field names match the struct tags
		user := UserProfile{
			ID:       "test",
			FullName: "Test User",
			Emoji:    "🧪",
		}

		jsonData, _ := json.Marshal(user)
		jsonString := string(jsonData)

		// Check that the JSON uses the correct field names as defined in struct tags
		if !containsSubstring(jsonString, `"id":"test"`) {
			t.Error("JSON should contain 'id' field")
		}
		if !containsSubstring(jsonString, `"fullName":"Test User"`) {
			t.Error("JSON should contain 'fullName' field")
		}
		if !containsSubstring(jsonString, `"emoji":"🧪"`) {
			t.Error("JSON should contain 'emoji' field")
		}
	})
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) != -1
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
