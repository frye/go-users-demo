package models

import (
	"encoding/json"
	"testing"
)

func TestUserProfileStruct(t *testing.T) {
	user := UserProfile{
		ID:       "1",
		FullName: "John Doe",
		Emoji:    "😀",
	}

	if user.ID != "1" {
		t.Errorf("Expected ID to be '1', got '%s'", user.ID)
	}

	if user.FullName != "John Doe" {
		t.Errorf("Expected FullName to be 'John Doe', got '%s'", user.FullName)
	}

	if user.Emoji != "😀" {
		t.Errorf("Expected Emoji to be '😀', got '%s'", user.Emoji)
	}
}

func TestUserProfileJSONMarshaling(t *testing.T) {
	user := UserProfile{
		ID:       "1",
		FullName: "John Doe",
		Emoji:    "😀",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal UserProfile: %v", err)
	}

	expected := `{"id":"1","fullName":"John Doe","emoji":"😀"}`
	if string(jsonData) != expected {
		t.Errorf("Expected JSON to be %s, got %s", expected, string(jsonData))
	}
}

func TestUserProfileJSONUnmarshaling(t *testing.T) {
	jsonData := `{"id":"2","fullName":"Jane Smith","emoji":"🚀"}`

	var user UserProfile
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal UserProfile: %v", err)
	}

	if user.ID != "2" {
		t.Errorf("Expected ID to be '2', got '%s'", user.ID)
	}

	if user.FullName != "Jane Smith" {
		t.Errorf("Expected FullName to be 'Jane Smith', got '%s'", user.FullName)
	}

	if user.Emoji != "🚀" {
		t.Errorf("Expected Emoji to be '🚀', got '%s'", user.Emoji)
	}
}
