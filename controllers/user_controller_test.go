package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"userprofile-api/models"
)

// setupTestRouter creates a new Gin router for testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

// resetUsers resets the users slice to a known state for testing
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

// TestHomePageHandler tests the HomePageHandler function
func TestHomePageHandler(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.LoadHTMLGlob("../templates/*")
	router.GET("/", HomePageHandler)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that the response contains HTML
	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type text/html; charset=utf-8, got %s", w.Header().Get("Content-Type"))
	}

	// Check that the response body contains some expected content
	body := w.Body.String()
	if len(body) == 0 {
		t.Error("Expected non-empty response body")
	}
}

// TestGetUsers tests the GetUsers function
func TestGetUsers(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var result []models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 users, got %d", len(result))
	}

	// Verify the first user
	if result[0].ID != "1" || result[0].FullName != "John Doe" {
		t.Errorf("Expected first user to be John Doe with ID 1, got %v", result[0])
	}
}

// TestGetUserSuccess tests the GetUser function with a valid ID
func TestGetUserSuccess(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, err := http.NewRequest("GET", "/api/v1/users/2", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var result models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.ID != "2" {
		t.Errorf("Expected user ID 2, got %s", result.ID)
	}

	if result.FullName != "Jane Smith" {
		t.Errorf("Expected user FullName 'Jane Smith', got %s", result.FullName)
	}

	if result.Emoji != "🚀" {
		t.Errorf("Expected user Emoji '🚀', got %s", result.Emoji)
	}
}

// TestGetUserNotFound tests the GetUser function with a non-existent ID
func TestGetUserNotFound(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, err := http.NewRequest("GET", "/api/v1/users/999", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got %s", result["error"])
	}
}

// TestCreateUserSuccess tests the CreateUser function with valid input
func TestCreateUserSuccess(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	newUser := models.UserProfile{
		ID:       "4",
		FullName: "Alice Johnson",
		Emoji:    "🌟",
	}

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var result models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.ID != "4" {
		t.Errorf("Expected user ID 4, got %s", result.ID)
	}

	if result.FullName != "Alice Johnson" {
		t.Errorf("Expected user FullName 'Alice Johnson', got %s", result.FullName)
	}

	if result.Emoji != "🌟" {
		t.Errorf("Expected user Emoji '🌟', got %s", result.Emoji)
	}

	// Verify the user was added to the users slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

// TestCreateUserInvalidJSON tests the CreateUser function with invalid JSON
func TestCreateUserInvalidJSON(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	invalidJSON := []byte(`{"id": "4", "fullName": "Alice Johnson", "emoji":`)

	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["error"] == "" {
		t.Error("Expected an error message, got empty string")
	}

	// Verify the users slice was not modified
	if len(users) != 3 {
		t.Errorf("Expected 3 users after invalid creation, got %d", len(users))
	}
}

// TestUpdateUserSuccess tests the UpdateUser function with valid input
func TestUpdateUserSuccess(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		ID:       "2",
		FullName: "Jane Doe",
		Emoji:    "✨",
	}

	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/v1/users/2", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var result models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result.ID != "2" {
		t.Errorf("Expected user ID 2, got %s", result.ID)
	}

	if result.FullName != "Jane Doe" {
		t.Errorf("Expected user FullName 'Jane Doe', got %s", result.FullName)
	}

	if result.Emoji != "✨" {
		t.Errorf("Expected user Emoji '✨', got %s", result.Emoji)
	}

	// Verify the user was updated in the users slice
	found := false
	for _, user := range users {
		if user.ID == "2" {
			if user.FullName != "Jane Doe" || user.Emoji != "✨" {
				t.Errorf("User was not properly updated in the users slice")
			}
			found = true
			break
		}
	}

	if !found {
		t.Error("Updated user not found in users slice")
	}
}

// TestUpdateUserNotFound tests the UpdateUser function with a non-existent ID
func TestUpdateUserNotFound(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		ID:       "999",
		FullName: "Non Existent",
		Emoji:    "❌",
	}

	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got %s", result["error"])
	}
}

// TestUpdateUserInvalidJSON tests the UpdateUser function with invalid JSON
func TestUpdateUserInvalidJSON(t *testing.T) {
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	invalidJSON := []byte(`{"id": "2", "fullName": "Jane Doe", "emoji":`)

	req, err := http.NewRequest("PUT", "/api/v1/users/2", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var result map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["error"] == "" {
		t.Error("Expected an error message, got empty string")
	}

	// Verify the user was not modified
	for _, user := range users {
		if user.ID == "2" {
			if user.FullName != "Jane Smith" || user.Emoji != "🚀" {
				t.Errorf("User should not have been modified")
			}
			break
		}
	}
}
