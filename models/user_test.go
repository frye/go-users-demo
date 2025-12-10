package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserProfile_JSONMarshaling(t *testing.T) {
	// Test that UserProfile can be marshaled to JSON correctly
	user := UserProfile{
		ID:       "1",
		FullName: "John Doe",
		Emoji:    "😀",
	}

	jsonData, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.NotNil(t, jsonData)

	// Verify the JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	assert.NoError(t, err)
	assert.Equal(t, "1", result["id"])
	assert.Equal(t, "John Doe", result["fullName"])
	assert.Equal(t, "😀", result["emoji"])
}

func TestUserProfile_JSONUnmarshaling(t *testing.T) {
	// Test that JSON can be unmarshaled into UserProfile correctly
	jsonStr := `{"id":"2","fullName":"Jane Smith","emoji":"🚀"}`
	
	var user UserProfile
	err := json.Unmarshal([]byte(jsonStr), &user)
	assert.NoError(t, err)
	assert.Equal(t, "2", user.ID)
	assert.Equal(t, "Jane Smith", user.FullName)
	assert.Equal(t, "🚀", user.Emoji)
}

func TestUserProfile_Struct(t *testing.T) {
	// Test that UserProfile struct can be created and fields accessed
	user := UserProfile{
		ID:       "3",
		FullName: "Robert Johnson",
		Emoji:    "🎸",
	}

	assert.Equal(t, "3", user.ID)
	assert.Equal(t, "Robert Johnson", user.FullName)
	assert.Equal(t, "🎸", user.Emoji)
}
