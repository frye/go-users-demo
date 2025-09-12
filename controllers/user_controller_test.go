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

// setupTestRouter creates a test router with HTML template loading disabled for testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Load HTML templates for HomePageHandler test
	router.LoadHTMLGlob("../templates/*")
	
	return router
}

// resetUsers resets the users slice to initial test data
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

func TestHomePageHandler(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/", HomePageHandler)

	// Create request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check content type contains HTML
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type to contain text/html, got %s", contentType)
	}

	// Check that response body contains expected user data
	body := w.Body.String()
	if !contains(body, "John Doe") || !contains(body, "Jane Smith") || !contains(body, "Robert Johnson") {
		t.Error("Expected HTML body to contain user names")
	}
}

func TestGetUsers(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	// Create request
	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Parse response body
	var responseUsers []models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}

	// Check that all users are returned
	if len(responseUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(responseUsers))
	}

	// Check specific user data
	expectedUsers := []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}

	for i, expectedUser := range expectedUsers {
		if responseUsers[i] != expectedUser {
			t.Errorf("Expected user %+v, got %+v", expectedUser, responseUsers[i])
		}
	}
}

func TestGetUser_Success(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	// Test getting existing user
	req, err := http.NewRequest("GET", "/api/v1/users/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var user models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}

	// Check user data
	expectedUser := models.UserProfile{ID: "2", FullName: "Jane Smith", Emoji: "🚀"}
	if user != expectedUser {
		t.Errorf("Expected user %+v, got %+v", expectedUser, user)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	// Test getting non-existing user
	req, err := http.NewRequest("GET", "/api/v1/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Parse error response
	var errorResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Fatal("Failed to parse error response JSON:", err)
	}

	// Check error message
	if errorResponse["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", errorResponse["error"])
	}
}

func TestCreateUser_Success(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	// Create new user data
	newUser := models.UserProfile{
		ID:       "4",
		FullName: "Alice Cooper",
		Emoji:    "🎭",
	}

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		t.Fatal("Failed to marshal user data:", err)
	}

	// Create request
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse response
	var createdUser models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &createdUser)
	if err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}

	// Check created user data
	if createdUser != newUser {
		t.Errorf("Expected created user %+v, got %+v", newUser, createdUser)
	}

	// Verify user was added to the users slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	// Create request with invalid JSON
	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse error response
	var errorResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Fatal("Failed to parse error response JSON:", err)
	}

	// Check that error message exists
	if errorResponse["error"] == "" {
		t.Error("Expected error message for invalid JSON")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	// Create updated user data
	updatedUser := models.UserProfile{
		FullName: "John Smith",
		Emoji:    "😎",
	}

	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatal("Failed to marshal user data:", err)
	}

	// Create request
	req, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response
	var responseUser models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}

	// Check that ID is preserved and other fields are updated
	expectedUser := models.UserProfile{
		ID:       "1", // ID should be preserved
		FullName: "John Smith",
		Emoji:    "😎",
	}

	if responseUser != expectedUser {
		t.Errorf("Expected updated user %+v, got %+v", expectedUser, responseUser)
	}

	// Verify the user was actually updated in the slice
	found := false
	for _, user := range users {
		if user.ID == "1" && user.FullName == "John Smith" && user.Emoji == "😎" {
			found = true
			break
		}
	}
	if !found {
		t.Error("User was not actually updated in the users slice")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	// Create updated user data
	updatedUser := models.UserProfile{
		FullName: "Non Existent User",
		Emoji:    "❓",
	}

	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatal("Failed to marshal user data:", err)
	}

	// Create request for non-existent user
	req, err := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Parse error response
	var errorResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Fatal("Failed to parse error response JSON:", err)
	}

	// Check error message
	if errorResponse["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", errorResponse["error"])
	}
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	// Reset users to initial state
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	// Create request with invalid JSON
	req, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse error response
	var errorResponse map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	if err != nil {
		t.Fatal("Failed to parse error response JSON:", err)
	}

	// Check that error message exists
	if errorResponse["error"] == "" {
		t.Error("Expected error message for invalid JSON")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    s[:len(substr)] == substr || 
		    s[len(s)-len(substr):] == substr ||
		    containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}