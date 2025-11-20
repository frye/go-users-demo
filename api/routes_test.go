package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"userprofile-api/models"
)

func TestSetupRouter(t *testing.T) {
	router := SetupRouter()

	if router == nil {
		t.Fatal("Expected router to be initialized, got nil")
	}
}

func TestRouteGetUsers(t *testing.T) {
	router := SetupRouter()

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
}

func TestRouteGetUserByID(t *testing.T) {
	router := SetupRouter()

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
}

func TestRouteGetUserNotFound(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestHomePageRoute(t *testing.T) {
	router := SetupRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that response is HTML
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected Content-Type to be 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

func TestAPIVersionGrouping(t *testing.T) {
	router := SetupRouter()

	// Test that the v1 API group is properly configured
	routes := router.Routes()
	
	hasV1Routes := false
	for _, route := range routes {
		if route.Path == "/api/v1/users" || route.Path == "/api/v1/users/:id" {
			hasV1Routes = true
			break
		}
	}

	if !hasV1Routes {
		t.Error("Expected API v1 routes to be configured")
	}
}
