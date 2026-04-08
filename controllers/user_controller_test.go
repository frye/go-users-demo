package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"userprofile-api/database"
	"userprofile-api/models"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestRouter(t *testing.T) (*gin.Engine, *database.TestDB) {
	t.Helper()

	tdb, err := database.NewTestDB()
	if err != nil {
		t.Fatalf("failed to create test DB: %v", err)
	}

	uc := NewUserController(tdb.DB)

	router := gin.New()
	router.GET("/api/v1/users", uc.GetUsers)
	router.GET("/api/v1/users/:id", uc.GetUser)
	router.POST("/api/v1/users", uc.CreateUser)
	router.PUT("/api/v1/users/:id", uc.UpdateUser)
	router.DELETE("/api/v1/users/:id", uc.DeleteUser)

	return router, tdb
}

func seedUser(t *testing.T, tdb *database.TestDB, fullName, emoji string) models.UserProfile {
	t.Helper()
	user := models.UserProfile{FullName: fullName, Emoji: emoji}
	if err := tdb.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return user
}

func TestGetUsersEmpty(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var users []models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &users)
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetUsersWithData(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	seedUser(t, tdb, "Alice", "👩")
	seedUser(t, tdb, "Bob", "👨")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var users []models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &users)
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	// Verify ordering by ID
	if users[0].FullName != "Alice" || users[1].FullName != "Bob" {
		t.Error("users are not returned in ID order")
	}
}

func TestGetUserByID(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	user := seedUser(t, tdb, "Alice", "👩")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var got models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &got)
	if got.FullName != user.FullName {
		t.Errorf("expected %s, got %s", user.FullName, got.FullName)
	}
}

func TestGetUserNotFound(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestGetUserInvalidID(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/abc", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateUser(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	body, _ := json.Marshal(models.CreateUserRequest{
		FullName: "Charlie",
		Emoji:    "🎉",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var user models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &user)
	if user.FullName != "Charlie" {
		t.Errorf("expected FullName Charlie, got %s", user.FullName)
	}
	if user.ID == 0 {
		t.Error("expected non-zero ID")
	}

	// Verify it persisted
	var count int64
	tdb.DB.Model(&models.UserProfile{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 user in DB, got %d", count)
	}
}

func TestCreateUserInvalidBody(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	body := []byte(`{"fullName": ""}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	seedUser(t, tdb, "Alice", "👩")

	body, _ := json.Marshal(models.UpdateUserRequest{
		FullName: "Alice Updated",
		Emoji:    "🌟",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var user models.UserProfile
	json.Unmarshal(w.Body.Bytes(), &user)
	if user.FullName != "Alice Updated" {
		t.Errorf("expected Alice Updated, got %s", user.FullName)
	}
	if user.ID != 1 {
		t.Errorf("expected ID 1 to be preserved, got %d", user.ID)
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	body, _ := json.Marshal(models.UpdateUserRequest{
		FullName: "Ghost",
		Emoji:    "👻",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	seedUser(t, tdb, "Alice", "👩")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify deletion
	var count int64
	tdb.DB.Model(&models.UserProfile{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 users after delete, got %d", count)
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestDeleteUserInvalidID(t *testing.T) {
	router, tdb := setupTestRouter(t)
	defer tdb.Close()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/abc", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
