package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"userprofile-api/models"
)

// setupTestRouter creates a test router with gin in test mode
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// resetUsers resets the users slice to the original test data
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

func TestHomePageHandler(t *testing.T) {
	router := setupTestRouter()
	router.LoadHTMLGlob("../templates/*")
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check if the response contains HTML content
	if !strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		t.Error("Expected HTML content type")
	}
}

func TestGetUsers(t *testing.T) {
	resetUsers() // Ensure we start with clean test data

	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Test successful response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var responseUsers []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check if we get the expected number of users
	expectedCount := 3
	if len(responseUsers) != expectedCount {
		t.Errorf("Expected %d users, got %d", expectedCount, len(responseUsers))
	}

	// Verify first user data
	if responseUsers[0].ID != "1" || responseUsers[0].FullName != "John Doe" || responseUsers[0].Emoji != "😀" {
		t.Error("First user data doesn't match expected values")
	}
}

func TestGetUser(t *testing.T) {
	resetUsers() // Ensure we start with clean test data

	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	// Test successful case - user exists
	t.Run("UserExists", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var responseUser models.UserProfile
		err := json.Unmarshal(w.Body.Bytes(), &responseUser)
		if err != nil {
			t.Errorf("Failed to parse response JSON: %v", err)
		}

		// Verify user data
		if responseUser.ID != "1" || responseUser.FullName != "John Doe" || responseUser.Emoji != "😀" {
			t.Error("User data doesn't match expected values")
		}
	})

	// Test failure case - user not found
	t.Run("UserNotFound", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		if err != nil {
			t.Errorf("Failed to parse error response JSON: %v", err)
		}

		if errorResponse["error"] != "User not found" {
			t.Errorf("Expected error message 'User not found', got '%s'", errorResponse["error"])
		}
	})
}

func TestCreateUser(t *testing.T) {
	resetUsers() // Ensure we start with clean test data

	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	// Test successful user creation
	t.Run("ValidUser", func(t *testing.T) {
		newUser := models.UserProfile{
			ID:       "4",
			FullName: "Alice Brown",
			Emoji:    "🌟",
		}

		jsonData, _ := json.Marshal(newUser)
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var responseUser models.UserProfile
		err := json.Unmarshal(w.Body.Bytes(), &responseUser)
		if err != nil {
			t.Errorf("Failed to parse response JSON: %v", err)
		}

		// Verify created user data
		if responseUser.ID != "4" || responseUser.FullName != "Alice Brown" || responseUser.Emoji != "🌟" {
			t.Error("Created user data doesn't match expected values")
		}

		// Verify user was added to the slice
		if len(users) != 4 {
			t.Errorf("Expected 4 users after creation, got %d", len(users))
		}
	})

	// Test invalid JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJSON := `{"id": "5", "fullName": "Invalid User", "emoji":}`
		req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		if err != nil {
			t.Errorf("Failed to parse error response JSON: %v", err)
		}

		if errorResponse["error"] == "" {
			t.Error("Expected error message for invalid JSON")
		}
	})

	// Test empty request body
	t.Run("EmptyBody", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/users", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	resetUsers() // Ensure we start with clean test data

	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	// Test successful user update
	t.Run("ValidUpdate", func(t *testing.T) {
		updatedUser := models.UserProfile{
			ID:       "1", // This should be ignored and set to the URL parameter
			FullName: "John Updated",
			Emoji:    "🔥",
		}

		jsonData, _ := json.Marshal(updatedUser)
		req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var responseUser models.UserProfile
		err := json.Unmarshal(w.Body.Bytes(), &responseUser)
		if err != nil {
			t.Errorf("Failed to parse response JSON: %v", err)
		}

		// Verify updated user data (ID should remain "1")
		if responseUser.ID != "1" || responseUser.FullName != "John Updated" || responseUser.Emoji != "🔥" {
			t.Error("Updated user data doesn't match expected values")
		}

		// Verify user was updated in the slice
		if users[0].FullName != "John Updated" || users[0].Emoji != "🔥" {
			t.Error("User was not properly updated in the slice")
		}
	})

	// Test user not found
	t.Run("UserNotFound", func(t *testing.T) {
		updatedUser := models.UserProfile{
			FullName: "Non-existent User",
			Emoji:    "❌",
		}

		jsonData, _ := json.Marshal(updatedUser)
		req, _ := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		if err != nil {
			t.Errorf("Failed to parse error response JSON: %v", err)
		}

		if errorResponse["error"] != "User not found" {
			t.Errorf("Expected error message 'User not found', got '%s'", errorResponse["error"])
		}
	})

	// Test invalid JSON
	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJSON := `{"fullName": "Invalid User", "emoji":}`
		req, _ := http.NewRequest("PUT", "/api/v1/users/1", strings.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var errorResponse map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		if err != nil {
			t.Errorf("Failed to parse error response JSON: %v", err)
		}

		if errorResponse["error"] == "" {
			t.Error("Expected error message for invalid JSON")
		}
	})
}

// TestUserSliceIsolation ensures that modifications in one test don't affect others
func TestUserSliceIsolation(t *testing.T) {
	resetUsers()
	
	// Verify we start with 3 users
	if len(users) != 3 {
		t.Errorf("Expected 3 users initially, got %d", len(users))
	}
	
	// Modify the slice
	users = append(users, models.UserProfile{ID: "test", FullName: "Test User", Emoji: "🧪"})
	
	if len(users) != 4 {
		t.Errorf("Expected 4 users after modification, got %d", len(users))
	}
	
	// Reset and verify
	resetUsers()
	if len(users) != 3 {
		t.Errorf("Expected 3 users after reset, got %d", len(users))
	}
}
