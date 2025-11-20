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
	router := gin.New()
	return router
}

// resetUsers resets the users slice to its initial state for testing
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

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that response contains HTML content
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

// TestGetUsers tests the GetUsers function
func TestGetUsers(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 users, got %d", len(response))
	}

	// Verify the first user
	if response[0].ID != "1" || response[0].FullName != "John Doe" || response[0].Emoji != "😀" {
		t.Errorf("First user data doesn't match expected values")
	}
}

// TestGetUser_Success tests GetUser function when user is found
func TestGetUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != "2" || response.FullName != "Jane Smith" || response.Emoji != "🚀" {
		t.Errorf("User data doesn't match expected values. Got: %+v", response)
	}
}

// TestGetUser_NotFound tests GetUser function when user is not found
func TestGetUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

// TestCreateUser_Success tests CreateUser function with valid input
func TestCreateUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	newUser := models.UserProfile{
		ID:       "4",
		FullName: "Alice Cooper",
		Emoji:    "🎭",
	}

	jsonData, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != "4" || response.FullName != "Alice Cooper" || response.Emoji != "🎭" {
		t.Errorf("Created user data doesn't match expected values. Got: %+v", response)
	}

	// Verify the user was added to the slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

// TestCreateUser_InvalidJSON tests CreateUser function with invalid JSON
func TestCreateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	invalidJSON := []byte(`{"id": "5", "fullName": "Invalid User"`)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

// TestCreateUser_EmptyBody tests CreateUser function with empty request body
func TestCreateUser_EmptyBody(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateUser_Success tests UpdateUser function when user exists
func TestUpdateUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "John Smith Updated",
		Emoji:    "😎",
	}

	jsonData, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// ID should remain unchanged
	if response.ID != "1" {
		t.Errorf("Expected ID to remain '1', got '%s'", response.ID)
	}

	if response.FullName != "John Smith Updated" || response.Emoji != "😎" {
		t.Errorf("Updated user data doesn't match expected values. Got: %+v", response)
	}

	// Verify the user was updated in the slice
	if users[0].FullName != "John Smith Updated" || users[0].Emoji != "😎" {
		t.Errorf("User was not updated in the slice")
	}
}

// TestUpdateUser_NotFound tests UpdateUser function when user doesn't exist
func TestUpdateUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Non-existent User",
		Emoji:    "❓",
	}

	jsonData, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

// TestUpdateUser_InvalidJSON tests UpdateUser function with invalid JSON
func TestUpdateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	invalidJSON := []byte(`{"fullName": "Invalid User"`)
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

// TestUpdateUser_EmptyBody tests UpdateUser function with empty request body
func TestUpdateUser_EmptyBody(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
