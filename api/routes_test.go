package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Set Gin to test mode to reduce console output during tests
	gin.SetMode(gin.TestMode)
}

func TestSetupRouter(t *testing.T) {
	// Test that SetupRouter returns a non-nil router
	router := SetupRouter()
	assert.NotNil(t, router)
}

func TestRoutes_GetUsers(t *testing.T) {
	// Test that GET /api/v1/users route is registered
	router := SetupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w, req)
	
	// Should return 200 OK (not 404 Not Found)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoutes_GetUserByID(t *testing.T) {
	// Test that GET /api/v1/users/:id route is registered
	router := SetupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	router.ServeHTTP(w, req)
	
	// Should return 200 OK (not 404 Not Found)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoutes_CreateUser(t *testing.T) {
	// Test that POST /api/v1/users route is registered
	router := SetupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// Should return 400 Bad Request (not 404 Not Found)
	// because we're not sending valid JSON, but route exists
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoutes_UpdateUser(t *testing.T) {
	// Test that PUT /api/v1/users/:id route is registered
	router := SetupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", nil)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// Should return 400 Bad Request (not 404 Not Found)
	// because we're not sending valid JSON, but route exists
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoutes_HomePageHandler(t *testing.T) {
	// Test that GET / route is registered
	router := SetupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	
	// Should return 200 OK (HTML page)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
}

func TestRoutes_NotFoundHandler(t *testing.T) {
	// Test that unregistered routes return 404
	router := SetupRouter()
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
}
