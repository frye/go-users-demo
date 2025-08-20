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
	router := gin.New()
	return router
}

// resetUsers resets the global users slice to initial test data
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

// TestHomePageHandler tests the home page HTML rendering
func TestHomePageHandler(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	
	// Load templates for testing - we'll setup a minimal template loading
	router.LoadHTMLGlob("../templates/*")
	router.GET("/", HomePageHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if HTML content type is set
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}

	// Check if the response contains expected HTML elements
	body := w.Body.String()
	if !contains(body, "User Profiles") {
		t.Error("Expected HTML to contain 'User Profiles' title")
	}
	if !contains(body, "John Doe") {
		t.Error("Expected HTML to contain user 'John Doe'")
	}
}

// TestGetUsers tests retrieving all users
func TestGetUsers(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users", GetUsers)

	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Test status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Test response body
	var responseUsers []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &responseUsers)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(responseUsers) != 3 {
		t.Errorf("Expected 3 users, got %d", len(responseUsers))
	}

	// Test first user
	expectedUser := models.UserProfile{ID: "1", FullName: "John Doe", Emoji: "😀"}
	if responseUsers[0] != expectedUser {
		t.Errorf("Expected user %+v, got %+v", expectedUser, responseUsers[0])
	}
}

// TestGetUser tests retrieving a single user by ID
func TestGetUser(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.GET("/api/v1/users/:id", GetUser)

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		expectedUser   *models.UserProfile
		expectedError  string
	}{
		{
			name:           "Existing user",
			userID:         "1",
			expectedStatus: http.StatusOK,
			expectedUser:   &models.UserProfile{ID: "1", FullName: "John Doe", Emoji: "😀"},
		},
		{
			name:           "Another existing user",
			userID:         "3",
			expectedStatus: http.StatusOK,
			expectedUser:   &models.UserProfile{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
		},
		{
			name:           "Non-existing user",
			userID:         "999",
			expectedStatus: http.StatusNotFound,
			expectedError:  "User not found",
		},
		{
			name:           "Special character user ID",
			userID:         "user@123",
			expectedStatus: http.StatusNotFound,
			expectedError:  "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/api/v1/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedUser != nil {
				var responseUser models.UserProfile
				err := json.Unmarshal(w.Body.Bytes(), &responseUser)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if responseUser != *tt.expectedUser {
					t.Errorf("Expected user %+v, got %+v", *tt.expectedUser, responseUser)
				}
			}

			if tt.expectedError != "" {
				var errorResponse map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				if err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse["error"] != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, errorResponse["error"])
				}
			}
		})
	}
}

// TestCreateUser tests creating a new user
func TestCreateUser(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.POST("/api/v1/users", CreateUser)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		expectedUser   *models.UserProfile
		expectedError  bool
	}{
		{
			name: "Valid user creation",
			requestBody: models.UserProfile{
				ID:       "4",
				FullName: "Alice Cooper",
				Emoji:    "🎭",
			},
			expectedStatus: http.StatusCreated,
			expectedUser: &models.UserProfile{
				ID:       "4",
				FullName: "Alice Cooper",
				Emoji:    "🎭",
			},
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Missing fields",
			requestBody: map[string]interface{}{
				"id": "5",
				// missing fullName and emoji
			},
			expectedStatus: http.StatusCreated, // Gin binding doesn't fail on missing fields by default
			expectedUser: &models.UserProfile{
				ID:       "5",
				FullName: "",
				Emoji:    "",
			},
		},
		{
			name: "Empty request body",
			requestBody: map[string]interface{}{},
			expectedStatus: http.StatusCreated,
			expectedUser: &models.UserProfile{
				ID:       "",
				FullName: "",
				Emoji:    "",
			},
		},
		{
			name: "User with special characters",
			requestBody: models.UserProfile{
				ID:       "6",
				FullName: "José María Aznar-López",
				Emoji:    "🇪🇸",
			},
			expectedStatus: http.StatusCreated,
			expectedUser: &models.UserProfile{
				ID:       "6",
				FullName: "José María Aznar-López",
				Emoji:    "🇪🇸",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetUsers() // Reset for each test to maintain isolation

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			req, _ := http.NewRequest("POST", "/api/v1/users", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedUser != nil {
				var responseUser models.UserProfile
				err := json.Unmarshal(w.Body.Bytes(), &responseUser)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if responseUser != *tt.expectedUser {
					t.Errorf("Expected user %+v, got %+v", *tt.expectedUser, responseUser)
				}

				// Verify user was actually added to the slice
				found := false
				for _, user := range users {
					if user.ID == tt.expectedUser.ID {
						found = true
						break
					}
				}
				if !found {
					t.Error("User was not added to the users slice")
				}
			}

			if tt.expectedError {
				var errorResponse map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				if err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse["error"] == "" {
					t.Error("Expected error response, but got none")
				}
			}
		})
	}
}

// TestUpdateUser tests updating an existing user
func TestUpdateUser(t *testing.T) {
	resetUsers()
	router := setupTestRouter()
	router.PUT("/api/v1/users/:id", UpdateUser)

	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		expectedStatus int
		expectedUser   *models.UserProfile
		expectedError  bool
	}{
		{
			name:   "Valid user update",
			userID: "1",
			requestBody: models.UserProfile{
				FullName: "John Smith",
				Emoji:    "😎",
			},
			expectedStatus: http.StatusOK,
			expectedUser: &models.UserProfile{
				ID:       "1", // ID should remain the same
				FullName: "John Smith",
				Emoji:    "😎",
			},
		},
		{
			name:   "Update only full name",
			userID: "2",
			requestBody: models.UserProfile{
				FullName: "Jane Doe",
				Emoji:    "",
			},
			expectedStatus: http.StatusOK,
			expectedUser: &models.UserProfile{
				ID:       "2", // ID should remain the same
				FullName: "Jane Doe",
				Emoji:    "",
			},
		},
		{
			name:   "Update with special characters",
			userID: "3",
			requestBody: models.UserProfile{
				FullName: "Roberto José Martínez",
				Emoji:    "🇪🇸",
			},
			expectedStatus: http.StatusOK,
			expectedUser: &models.UserProfile{
				ID:       "3",
				FullName: "Roberto José Martínez",
				Emoji:    "🇪🇸",
			},
		},
		{
			name:   "Update non-existing user",
			userID: "999",
			requestBody: models.UserProfile{
				FullName: "Non Existing",
				Emoji:    "❌",
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "Invalid JSON",
			userID:         "1",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetUsers() // Reset for each test to maintain isolation

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.requestBody)
			}

			req, _ := http.NewRequest("PUT", "/api/v1/users/"+tt.userID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedUser != nil {
				var responseUser models.UserProfile
				err := json.Unmarshal(w.Body.Bytes(), &responseUser)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if responseUser != *tt.expectedUser {
					t.Errorf("Expected user %+v, got %+v", *tt.expectedUser, responseUser)
				}

				// Verify user was actually updated in the slice
				found := false
				for _, user := range users {
					if user.ID == tt.expectedUser.ID && user.FullName == tt.expectedUser.FullName && user.Emoji == tt.expectedUser.Emoji {
						found = true
						break
					}
				}
				if !found {
					t.Error("User was not updated in the users slice")
				}
			}

			if tt.expectedError {
				var errorResponse map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				if err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse["error"] == "" {
					t.Error("Expected error response, but got none")
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}