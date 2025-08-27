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

// setupTestGin creates a gin engine for testing
func setupTestGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// setupTestUsers resets the users slice to a known state for testing
func setupTestUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

// resetUsers restores the original users slice after test
func resetUsers() {
	setupTestUsers() // Reset to original state
}

func TestHomePageHandler(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.LoadHTMLGlob("../templates/*")
	router.GET("/", HomePageHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that it returns HTML content
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	// Check that response contains user data
	body := w.Body.String()
	if !strings.Contains(body, "John Doe") {
		t.Error("Expected response to contain 'John Doe'")
	}
	if !strings.Contains(body, "Jane Smith") {
		t.Error("Expected response to contain 'Jane Smith'")
	}
	if !strings.Contains(body, "Robert Johnson") {
		t.Error("Expected response to contain 'Robert Johnson'")
	}
}

func TestGetUsers(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.GET("/api/v1/users", GetUsers)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check Content-Type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", contentType)
	}

	// Parse response body
	var responseUsers []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check that we get all users
	if len(responseUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(responseUsers))
	}

	// Check specific users are present
	expectedUsers := map[string]models.UserProfile{
		"1": {ID: "1", FullName: "John Doe", Emoji: "😀"},
		"2": {ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		"3": {ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}

	for _, user := range responseUsers {
		expected, exists := expectedUsers[user.ID]
		if !exists {
			t.Errorf("Unexpected user with ID %s", user.ID)
			continue
		}
		if user != expected {
			t.Errorf("User %s: expected %+v, got %+v", user.ID, expected, user)
		}
	}
}

func TestGetUser_Success(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.GET("/api/v1/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check user data
	expectedUser := models.UserProfile{ID: "1", FullName: "John Doe", Emoji: "😀"}
	if responseUser != expectedUser {
		t.Errorf("Expected user %+v, got %+v", expectedUser, responseUser)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.GET("/api/v1/users/:id", GetUser)

	req := httptest.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Parse response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check error message
	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}
}

func TestCreateUser_Success(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.POST("/api/v1/users", CreateUser)

	newUser := models.UserProfile{ID: "4", FullName: "Alice Cooper", Emoji: "🎭"}
	jsonData, _ := json.Marshal(newUser)

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Parse response body
	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check that the created user is returned
	if responseUser != newUser {
		t.Errorf("Expected user %+v, got %+v", newUser, responseUser)
	}

	// Check that user was actually added to the slice
	if len(users) != 4 {
		t.Errorf("Expected 4 users after creation, got %d", len(users))
	}

	// Verify the new user is in the slice
	found := false
	for _, user := range users {
		if user.ID == "4" && user.FullName == "Alice Cooper" && user.Emoji == "🎭" {
			found = true
			break
		}
	}
	if !found {
		t.Error("New user was not found in users slice")
	}
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.POST("/api/v1/users", CreateUser)

	// Send invalid JSON
	invalidJSON := `{"id": "4", "fullName": "Alice Cooper"`  // Missing closing brace and emoji field

	req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check that error field exists
	if _, exists := response["error"]; !exists {
		t.Error("Expected error field in response")
	}

	// Check that users slice wasn't modified
	if len(users) != 3 {
		t.Errorf("Expected users slice to remain unchanged with 3 users, got %d", len(users))
	}
}

func TestUpdateUser_Success(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{FullName: "John Smith", Emoji: "😎"}
	jsonData, _ := json.Marshal(updatedUser)

	req := httptest.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Parse response body
	var responseUser models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUser)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check that the updated user is returned with correct ID
	expectedUser := models.UserProfile{ID: "1", FullName: "John Smith", Emoji: "😎"}
	if responseUser != expectedUser {
		t.Errorf("Expected user %+v, got %+v", expectedUser, responseUser)
	}

	// Check that user was actually updated in the slice
	found := false
	for _, user := range users {
		if user.ID == "1" {
			if user.FullName == "John Smith" && user.Emoji == "😎" {
				found = true
			}
			break
		}
	}
	if !found {
		t.Error("User was not properly updated in users slice")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{FullName: "Non Existent", Emoji: "🤷"}
	jsonData, _ := json.Marshal(updatedUser)

	req := httptest.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	// Parse response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check error message
	if response["error"] != "User not found" {
		t.Errorf("Expected error message 'User not found', got '%s'", response["error"])
	}

	// Check that users slice wasn't modified
	if len(users) != 3 {
		t.Errorf("Expected users slice to remain unchanged with 3 users, got %d", len(users))
	}
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	setupTestUsers()
	defer resetUsers()

	router := setupTestGin()
	router.PUT("/api/v1/users/:id", UpdateUser)

	// Send invalid JSON
	invalidJSON := `{"fullName": "John Smith"`  // Missing closing brace and emoji field

	req := httptest.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check that error field exists
	if _, exists := response["error"]; !exists {
		t.Error("Expected error field in response")
	}

	// Check that the original user wasn't modified
	for _, user := range users {
		if user.ID == "1" {
			if user.FullName != "John Doe" || user.Emoji != "😀" {
				t.Error("Original user should not have been modified due to invalid JSON")
			}
			break
		}
	}
}

