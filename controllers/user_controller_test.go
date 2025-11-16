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

// setupTestRouter creates a Gin router for testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	// Load templates for HomePageHandler test
	router.LoadHTMLGlob("../templates/*")
	
	return router
}

// resetUsers resets the users slice to the original state
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

	// Check if the response contains HTML content
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	// Verify that the response body contains user data
	body := w.Body.String()
	if len(body) == 0 {
		t.Error("Expected non-empty HTML response body")
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

	var responseUsers []models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(responseUsers))
	}

	if responseUsers[0].ID != "1" || responseUsers[0].FullName != "John Doe" {
		t.Errorf("User data mismatch for first user")
	}
}

// TestGetUser_Success tests retrieving an existing user
func TestGetUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, err := http.NewRequest("GET", "/api/v1/users/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var user models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.ID != "1" {
		t.Errorf("Expected user ID '1', got '%s'", user.ID)
	}
	if user.FullName != "John Doe" {
		t.Errorf("Expected full name 'John Doe', got '%s'", user.FullName)
	}
	if user.Emoji != "😀" {
		t.Errorf("Expected emoji '😀', got '%s'", user.Emoji)
	}
}

// TestGetUser_NotFound tests retrieving a non-existent user
func TestGetUser_NotFound(t *testing.T) {
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

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

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
		FullName: "Alice Johnson",
		Emoji:    "🎨",
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

	var createdUser models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &createdUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdUser.ID != "4" {
		t.Errorf("Expected user ID '4', got '%s'", createdUser.ID)
	}
	if createdUser.FullName != "Alice Johnson" {
		t.Errorf("Expected full name 'Alice Johnson', got '%s'", createdUser.FullName)
	}
}

// TestCreateUser_InvalidJSON tests creating a user with invalid JSON
func TestCreateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	invalidJSON := []byte(`{"id": "4", "fullName": "Alice", "emoji":}`)

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

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Error("Expected error message in response")
	}
}

// TestCreateUser_EmptyBody tests creating a user with empty request body
func TestCreateUser_EmptyBody(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestUpdateUser_Success tests updating an existing user successfully
func TestUpdateUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		ID:       "1",
		FullName: "John Updated",
		Emoji:    "🌟",
	}

	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	req, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUser models.UserProfile
	err = json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseUser.FullName != "John Updated" {
		t.Errorf("Expected full name 'John Updated', got '%s'", responseUser.FullName)
	}
	if responseUser.Emoji != "🌟" {
		t.Errorf("Expected emoji '🌟', got '%s'", responseUser.Emoji)
	}
	// Verify ID was preserved
	if responseUser.ID != "1" {
		t.Errorf("Expected user ID '1' to be preserved, got '%s'", responseUser.ID)
	}
}

// TestUpdateUser_NotFound tests updating a non-existent user
func TestUpdateUser_NotFound(t *testing.T) {
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

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

// TestUpdateUser_InvalidJSON tests updating a user with invalid JSON
func TestUpdateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	invalidJSON := []byte(`{"fullName": "John", "emoji":}`)

	req, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Error("Expected error message in response")
	}
}

// TestUpdateUser_EmptyBody tests updating a user with empty request body
func TestUpdateUser_EmptyBody(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	req, err := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
