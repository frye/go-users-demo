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

// setupRouter creates a test router
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

// resetUsers resets the users slice to the original state for testing
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
	router := setupRouter()
	router.GET("/api/v1/users", GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUsers []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if len(responseUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(responseUsers))
	}

	if responseUsers[0].ID != "1" || responseUsers[0].FullName != "John Doe" {
		t.Errorf("Expected first user to be John Doe with ID 1, got %+v", responseUsers[0])
	}
}

// TestGetUser_Success tests getting an existing user by ID
func TestGetUser_Success(t *testing.T) {
	resetUsers()
	router := setupRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var user models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if user.ID != "2" || user.FullName != "Jane Smith" || user.Emoji != "🚀" {
		t.Errorf("Expected Jane Smith with ID 2 and emoji 🚀, got %+v", user)
	}
}

// TestGetUser_NotFound tests getting a non-existent user
func TestGetUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupRouter()
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
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

// TestCreateUser_Success tests creating a new user
func TestCreateUser_Success(t *testing.T) {
	resetUsers()
	router := setupRouter()
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

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var createdUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if createdUser.ID != "4" || createdUser.FullName != "Alice Wonder" || createdUser.Emoji != "🌟" {
		t.Errorf("Expected Alice Wonder with ID 4 and emoji 🌟, got %+v", createdUser)
	}

	// Verify user was added to the list
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

// TestCreateUser_InvalidJSON tests creating a user with invalid JSON
func TestCreateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupRouter()
	router.POST("/api/v1/users", CreateUser)

	invalidJSON := []byte(`{"id": "4", "fullName": "Invalid User", "emoji":}`)
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
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if response["error"] == "" {
		t.Errorf("Expected error message for invalid JSON, got empty string")
	}
}

// TestUpdateUser_Success tests updating an existing user
func TestUpdateUser_Success(t *testing.T) {
	resetUsers()
	router := setupRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Johnny Doe",
		Emoji:    "🎉",
	}

	jsonData, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if responseUser.ID != "1" {
		t.Errorf("Expected user ID to remain '1', got '%s'", responseUser.ID)
	}

	if responseUser.FullName != "Johnny Doe" || responseUser.Emoji != "🎉" {
		t.Errorf("Expected updated user Johnny Doe with emoji 🎉, got %+v", responseUser)
	}

	// Verify the user in the list was actually updated
	if users[0].FullName != "Johnny Doe" || users[0].Emoji != "🎉" {
		t.Errorf("Expected user in list to be updated, got %+v", users[0])
	}
}

// TestUpdateUser_NotFound tests updating a non-existent user
func TestUpdateUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Non Existent",
		Emoji:    "👻",
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
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

// TestUpdateUser_InvalidJSON tests updating a user with invalid JSON
func TestUpdateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	invalidJSON := []byte(`{"fullName": "Invalid", "emoji":}`)
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
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if response["error"] == "" {
		t.Errorf("Expected error message for invalid JSON, got empty string")
	}
}

// TestHomePageHandler tests the home page HTML handler
func TestHomePageHandler(t *testing.T) {
	resetUsers()
	router := setupRouter()
	
	// Load HTML templates for testing
	router.LoadHTMLGlob("../templates/*")
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if the response contains HTML content
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	// Verify that the response body contains expected user data
	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty response body for HTML page")
	}
}
