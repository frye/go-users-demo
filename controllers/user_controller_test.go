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

// setupTestRouter creates a test router with Gin in test mode
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Setup HTML templates for testing
	router.LoadHTMLGlob("../templates/*")
	
	// Setup routes
	router.GET("/", HomePageHandler)
	
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", GetUsers)
			users.GET("/:id", GetUser)
			users.POST("", CreateUser)
			users.PUT("/:id", UpdateUser)
		}
	}
	
	return router
}

// resetTestData resets the users slice to initial test data
func resetTestData() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

func TestHomePageHandler(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if the response contains HTML content
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	// Check if response body contains expected user data
	body := w.Body.String()
	if !contains(body, "John Doe") || !contains(body, "Jane Smith") || !contains(body, "Robert Johnson") {
		t.Error("Response body should contain user names")
	}
}

func TestGetUsers(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUsers []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(responseUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(responseUsers))
	}

	// Check specific user data
	expectedUsers := []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}

	for i, expected := range expectedUsers {
		if responseUsers[i] != expected {
			t.Errorf("User %d mismatch. Expected: %+v, Got: %+v", i, expected, responseUsers[i])
		}
	}
}

func TestGetUser_Success(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	expectedUser := models.UserProfile{ID: "1", FullName: "John Doe", Emoji: "😀"}
	if responseUser != expectedUser {
		t.Errorf("Expected user: %+v, Got: %+v", expectedUser, responseUser)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

func TestCreateUser_Success(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	newUser := models.UserProfile{
		ID:       "4",
		FullName: "Alice Cooper",
		Emoji:    "🎭",
	}

	jsonData, _ := json.Marshal(newUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if responseUser != newUser {
		t.Errorf("Expected user: %+v, Got: %+v", newUser, responseUser)
	}

	// Verify user was added to the slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

func TestCreateUser_BadRequest(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	// Send invalid JSON
	invalidJSON := `{"id": "4", "fullName": "Alice Cooper", "emoji": 123}` // emoji should be string, not number
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Error("Expected error message in response")
	}
}

func TestCreateUser_EmptyBody(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	updatedUser := models.UserProfile{
		FullName: "John Smith",
		Emoji:    "😎",
	}

	jsonData, _ := json.Marshal(updatedUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// The ID should be preserved and set to "1"
	expectedUser := models.UserProfile{ID: "1", FullName: "John Smith", Emoji: "😎"}
	if responseUser != expectedUser {
		t.Errorf("Expected user: %+v, Got: %+v", expectedUser, responseUser)
	}

	// Verify user was updated in the slice
	if users[0].FullName != "John Smith" || users[0].Emoji != "😎" {
		t.Error("User was not properly updated in the slice")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	updatedUser := models.UserProfile{
		FullName: "Non Existent",
		Emoji:    "🤷",
	}

	jsonData, _ := json.Marshal(updatedUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

func TestUpdateUser_BadRequest(t *testing.T) {
	resetTestData()
	router := setupTestRouter()

	// Send invalid JSON
	invalidJSON := `{"fullName": "John Smith", "emoji": 123}` // emoji should be string, not number
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["error"] == "" {
		t.Error("Expected error message in response")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && containsAt(s, substr, 0)))
}

func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[start+i] != substr[i] {
			if start+1 < len(s) {
				return containsAt(s, substr, start+1)
			}
			return false
		}
	}
	return true
}