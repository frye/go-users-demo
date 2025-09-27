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

// setupTestRouter creates a test router for testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// resetUsers resets the global users slice to its initial state before each test
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

func TestHomePageHandler(t *testing.T) {
	// Reset users before test
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

	if w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type to be text/html; charset=utf-8, got %s", w.Header().Get("Content-Type"))
	}
}

func TestGetUsers(t *testing.T) {
	// Reset users before test
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

	// Check Content-Type
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type to be application/json; charset=utf-8, got %s", w.Header().Get("Content-Type"))
	}

	// Parse response body
	var responseUsers []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check if we got the expected number of users
	expectedCount := 3
	if len(responseUsers) != expectedCount {
		t.Errorf("Expected %d users, got %d", expectedCount, len(responseUsers))
	}

	// Check if the first user matches expected data
	if responseUsers[0].ID != "1" || responseUsers[0].FullName != "John Doe" || responseUsers[0].Emoji != "😀" {
		t.Errorf("First user data doesn't match expected values")
	}
}

func TestGetUser_Success(t *testing.T) {
	// Reset users before test
	resetUsers()
	
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var user models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check user data
	if user.ID != "1" || user.FullName != "John Doe" || user.Emoji != "😀" {
		t.Errorf("User data doesn't match expected values: got %+v", user)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	// Reset users before test
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

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check error message
	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got %v", response["error"])
	}
}

func TestCreateUser_Success(t *testing.T) {
	// Reset users before test
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

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse response body
	var createdUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check created user data
	if createdUser.ID != "4" || createdUser.FullName != "Alice Cooper" || createdUser.Emoji != "🎭" {
		t.Errorf("Created user data doesn't match expected values: got %+v", createdUser)
	}

	// Verify user was actually added to the slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	// Reset users before test
	resetUsers()
	
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	// Send invalid JSON (missing value after emoji field)
	invalidJSON := `{"id": "4", "fullName": "Alice Cooper", "emoji": }`
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check that an error field exists
	if response["error"] == nil {
		t.Errorf("Expected error field in response")
	}

	// Verify no user was added
	if len(users) != 3 {
		t.Errorf("Expected 3 users after failed creation, got %d", len(users))
	}
}

func TestUpdateUser_Success(t *testing.T) {
	// Reset users before test
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "John Smith",
		Emoji:    "😎",
	}

	jsonData, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check updated user data
	if responseUser.ID != "1" || responseUser.FullName != "John Smith" || responseUser.Emoji != "😎" {
		t.Errorf("Updated user data doesn't match expected values: got %+v", responseUser)
	}

	// Verify user was actually updated in the slice
	for _, user := range users {
		if user.ID == "1" {
			if user.FullName != "John Smith" || user.Emoji != "😎" {
				t.Errorf("User was not properly updated in the slice")
			}
			break
		}
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	// Reset users before test
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Nonexistent User",
		Emoji:    "❓",
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

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check error message
	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got %v", response["error"])
	}
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	// Reset users before test
	resetUsers()
	
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	// Send invalid JSON (missing value after emoji field)
	invalidJSON := `{"fullName": "John Smith", "emoji": }`
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response JSON: %v", err)
	}

	// Check that an error field exists
	if response["error"] == nil {
		t.Errorf("Expected error field in response")
	}

	// Verify original user data wasn't changed
	for _, user := range users {
		if user.ID == "1" {
			if user.FullName != "John Doe" || user.Emoji != "😀" {
				t.Errorf("User should not have been modified due to invalid JSON")
			}
			break
		}
	}
}