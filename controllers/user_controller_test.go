package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"userprofile-api/models"
)

// setupTestRouter creates a test router with HTML templates
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Get the absolute path to the templates directory
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(b))
	templatesPath := filepath.Join(basePath, "templates/*")
	
	// Load HTML templates for HomePageHandler test
	router.LoadHTMLGlob(templatesPath)
	
	return router
}

// resetUsers resets the users slice to its initial state before each test
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

func TestHomePageHandler(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/", HomePageHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that the response contains HTML content
	body := w.Body.String()
	if !strings.Contains(body, "User Profiles") {
		t.Error("Expected HTML response to contain 'User Profiles' title")
	}

	// Check that user data is rendered
	if !strings.Contains(body, "John Doe") {
		t.Error("Expected HTML response to contain user data")
	}
}

func TestGetUsers(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUsers []models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &responseUsers); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(responseUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(responseUsers))
	}

	// Check that all expected users are present
	expectedIDs := []string{"1", "2", "3"}
	for i, user := range responseUsers {
		if user.ID != expectedIDs[i] {
			t.Errorf("Expected user ID %s, got %s", expectedIDs[i], user.ID)
		}
	}
}

func TestGetUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUser models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseUser.ID != "1" {
		t.Errorf("Expected user ID '1', got '%s'", responseUser.ID)
	}
	if responseUser.FullName != "John Doe" {
		t.Errorf("Expected full name 'John Doe', got '%s'", responseUser.FullName)
	}
	if responseUser.Emoji != "😀" {
		t.Errorf("Expected emoji '😀', got '%s'", responseUser.Emoji)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", errorResponse["error"])
	}
}

func TestCreateUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	newUser := models.UserProfile{
		ID:       "4",
		FullName: "Alice Cooper",
		Emoji:    "🎭",
	}

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var responseUser models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseUser.ID != "4" {
		t.Errorf("Expected user ID '4', got '%s'", responseUser.ID)
	}
	if responseUser.FullName != "Alice Cooper" {
		t.Errorf("Expected full name 'Alice Cooper', got '%s'", responseUser.FullName)
	}

	// Verify user was actually added to the slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	invalidJSON := `{"id": "4", "fullName": "Alice Cooper", "emoji":` // Incomplete JSON

	req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse["error"] == "" {
		t.Error("Expected error message for invalid JSON")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Johnny Doe",
		Emoji:    "😎",
	}

	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	req := httptest.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUser models.UserProfile
	if err := json.Unmarshal(w.Body.Bytes(), &responseUser); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseUser.ID != "1" {
		t.Errorf("Expected user ID to remain '1', got '%s'", responseUser.ID)
	}
	if responseUser.FullName != "Johnny Doe" {
		t.Errorf("Expected updated full name 'Johnny Doe', got '%s'", responseUser.FullName)
	}
	if responseUser.Emoji != "😎" {
		t.Errorf("Expected updated emoji '😎', got '%s'", responseUser.Emoji)
	}

	// Verify the user was actually updated in the slice
	for _, user := range users {
		if user.ID == "1" {
			if user.FullName != "Johnny Doe" {
				t.Error("User was not actually updated in the users slice")
			}
			break
		}
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Non Existent User",
		Emoji:    "👻",
	}

	jsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	req := httptest.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", errorResponse["error"])
	}
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	invalidJSON := `{"fullName": "Johnny Doe", "emoji":` // Incomplete JSON

	req := httptest.NewRequest("PUT", "/api/v1/users/1", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var errorResponse map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if errorResponse["error"] == "" {
		t.Error("Expected error message for invalid JSON")
	}
}