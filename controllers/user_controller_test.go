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

// This file contains comprehensive unit tests for all controller functions in user_controller.go.
// Tests cover both successful operations and error scenarios including:
// - Invalid JSON input handling
// - Not found scenarios
// - Proper HTTP status codes
// - Response data validation

// setupTestRouter creates a test router with test mode enabled
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// resetUsers resets the users slice to the initial test data
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

// TestGetUsers tests the GetUsers handler
func TestGetUsers(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var response []models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check number of users
	if len(response) != 3 {
		t.Errorf("Expected 3 users, got %d", len(response))
	}

	// Verify first user
	if response[0].ID != "1" || response[0].FullName != "John Doe" || response[0].Emoji != "😀" {
		t.Errorf("First user data mismatch: %+v", response[0])
	}
}

// TestGetUser_Success tests getting a single user by ID successfully
func TestGetUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var response models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify user data
	if response.ID != "2" || response.FullName != "Jane Smith" || response.Emoji != "🚀" {
		t.Errorf("User data mismatch: %+v", response)
	}
}

// TestGetUser_NotFound tests getting a user that doesn't exist
func TestGetUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Parse error response
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify error message
	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

// TestCreateUser_Success tests creating a new user successfully
func TestCreateUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	newUser := models.UserProfile{
		ID:       "4",
		FullName: "Alice Wonder",
		Emoji:    "🌟",
	}

	jsonData, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse response
	var response models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify created user data
	if response.ID != "4" || response.FullName != "Alice Wonder" || response.Emoji != "🌟" {
		t.Errorf("Created user data mismatch: %+v", response)
	}

	// Verify user was added to the slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

// TestCreateUser_InvalidJSON tests creating a user with invalid JSON
func TestCreateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	invalidJSON := []byte(`{"id": "4", "fullName": "Alice", "emoji":}`)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse error response
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify error exists
	if response["error"] == "" {
		t.Error("Expected error message, got empty string")
	}
}

// TestUpdateUser_Success tests updating an existing user successfully
func TestUpdateUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		ID:       "2", // This will be overridden by the handler
		FullName: "Jane Doe Smith",
		Emoji:    "🌈",
	}

	jsonData, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/api/v1/users/2", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var response models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify updated user data
	if response.ID != "2" || response.FullName != "Jane Doe Smith" || response.Emoji != "🌈" {
		t.Errorf("Updated user data mismatch: %+v", response)
	}

	// Verify user was actually updated in the slice
	if users[1].FullName != "Jane Doe Smith" || users[1].Emoji != "🌈" {
		t.Errorf("User not updated in slice: %+v", users[1])
	}
}

// TestUpdateUser_NotFound tests updating a user that doesn't exist
func TestUpdateUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		ID:       "999",
		FullName: "Non Existent",
		Emoji:    "❌",
	}

	jsonData, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Parse error response
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify error message
	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

// TestUpdateUser_InvalidJSON tests updating a user with invalid JSON
func TestUpdateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	invalidJSON := []byte(`{"fullName": "Test", "emoji":}`)
	req, _ := http.NewRequest("PUT", "/api/v1/users/2", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse error response
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify error exists
	if response["error"] == "" {
		t.Error("Expected error message, got empty string")
	}
}

// TestHomePageHandler tests the home page handler
func TestHomePageHandler(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	
	// Load the HTML template for testing
	router.LoadHTMLGlob("/home/runner/work/go-users-demo/go-users-demo/templates/*")
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	// Check that response contains some expected content
	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty HTML response")
	}
}
