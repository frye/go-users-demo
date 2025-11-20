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

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestGetUsers(t *testing.T) {
	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var users []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &users)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(users) < 1 {
		t.Errorf("Expected at least 1 user, got %d", len(users))
	}
}

func TestGetUserSuccess(t *testing.T) {
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var user models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if user.ID != "1" {
		t.Errorf("Expected user ID to be '1', got '%s'", user.ID)
	}

	if user.FullName == "" {
		t.Error("Expected user to have a FullName")
	}

	if user.Emoji == "" {
		t.Error("Expected user to have an Emoji")
	}
}

func TestGetUserNotFound(t *testing.T) {
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

func TestCreateUserSuccess(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	newUser := models.UserProfile{
		ID:       "100",
		FullName: "Test User",
		Emoji:    "🧪",
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
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdUser.ID != newUser.ID {
		t.Errorf("Expected user ID to be '%s', got '%s'", newUser.ID, createdUser.ID)
	}

	if createdUser.FullName != newUser.FullName {
		t.Errorf("Expected user FullName to be '%s', got '%s'", newUser.FullName, createdUser.FullName)
	}

	if createdUser.Emoji != newUser.Emoji {
		t.Errorf("Expected user Emoji to be '%s', got '%s'", newUser.Emoji, createdUser.Emoji)
	}
}

func TestCreateUserBadRequest(t *testing.T) {
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	invalidJSON := `{"id": "101", "fullName": }`
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Updated Name",
		Emoji:    "✨",
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
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseUser.ID != "1" {
		t.Errorf("Expected user ID to remain '1', got '%s'", responseUser.ID)
	}

	if responseUser.FullName != updatedUser.FullName {
		t.Errorf("Expected user FullName to be '%s', got '%s'", updatedUser.FullName, responseUser.FullName)
	}

	if responseUser.Emoji != updatedUser.Emoji {
		t.Errorf("Expected user Emoji to be '%s', got '%s'", updatedUser.Emoji, responseUser.Emoji)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	updatedUser := models.UserProfile{
		FullName: "Updated Name",
		Emoji:    "✨",
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

func TestUpdateUserBadRequest(t *testing.T) {
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	invalidJSON := `{"fullName": "Test", "emoji": }`
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
