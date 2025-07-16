package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	router := SetupRouter()

	// Test that router is created successfully
	if router == nil {
		t.Fatal("Router should not be nil")
	}

	// Test GET / route exists
	t.Run("HomePageRoute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should get 200 OK (even if template loading fails in test, the route exists)
		if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 200 or 500, got %d", w.Code)
		}
	})

	// Test GET /api/v1/users route exists
	t.Run("GetUsersRoute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	// Test GET /api/v1/users/:id route exists
	t.Run("GetUserRoute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	// Test POST /api/v1/users route exists (should fail with bad request due to empty body)
	t.Run("CreateUserRoute", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/v1/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should get 400 Bad Request due to empty body
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test PUT /api/v1/users/:id route exists (should fail with bad request due to empty body)
	t.Run("UpdateUserRoute", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/api/v1/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should get 400 Bad Request due to empty body
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test that non-existent route returns 404
	t.Run("NonExistentRoute", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
