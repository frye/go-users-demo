package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"userprofile-api/models"
)

func init() {
	// Set Gin to test mode to reduce console output during tests
	gin.SetMode(gin.TestMode)
}

// resetUsers resets the users slice to the default test data
func resetUsers() {
	users = []models.UserProfile{
		{ID: "1", FullName: "John Doe", Emoji: "😀"},
		{ID: "2", FullName: "Jane Smith", Emoji: "🚀"},
		{ID: "3", FullName: "Robert Johnson", Emoji: "🎸"},
	}
}

func TestGetUsers(t *testing.T) {
	resetUsers()

	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Execute
	GetUsers(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 3)
	assert.Equal(t, "1", response[0].ID)
	assert.Equal(t, "John Doe", response[0].FullName)
	assert.Equal(t, "😀", response[0].Emoji)
}

func TestGetUser_Found(t *testing.T) {
	resetUsers()

	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "2"}}

	// Execute
	GetUser(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "2", response.ID)
	assert.Equal(t, "Jane Smith", response.FullName)
	assert.Equal(t, "🚀", response.Emoji)
}

func TestGetUser_NotFound(t *testing.T) {
	resetUsers()

	// Setup
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}

	// Execute
	GetUser(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

func TestCreateUser_Success(t *testing.T) {
	resetUsers()

	// Setup
	newUser := models.UserProfile{
		ID:       "4",
		FullName: "Alice Cooper",
		Emoji:    "🎭",
	}
	jsonData, _ := json.Marshal(newUser)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	CreateUser(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "4", response.ID)
	assert.Equal(t, "Alice Cooper", response.FullName)
	assert.Equal(t, "🎭", response.Emoji)

	// Verify user was added to the slice
	assert.Len(t, users, 4)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	resetUsers()

	// Setup
	invalidJSON := []byte(`{"id": "invalid JSON`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	CreateUser(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["error"])
}

func TestUpdateUser_Success(t *testing.T) {
	resetUsers()

	// Setup
	updatedUser := models.UserProfile{
		FullName: "John Smith",
		Emoji:    "😎",
	}
	jsonData, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	UpdateUser(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.UserProfile
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "1", response.ID) // ID should remain the same
	assert.Equal(t, "John Smith", response.FullName)
	assert.Equal(t, "😎", response.Emoji)

	// Verify user was updated in the slice
	assert.Equal(t, "John Smith", users[0].FullName)
	assert.Equal(t, "😎", users[0].Emoji)
}

func TestUpdateUser_NotFound(t *testing.T) {
	resetUsers()

	// Setup
	updatedUser := models.UserProfile{
		FullName: "Nobody",
		Emoji:    "👻",
	}
	jsonData, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	UpdateUser(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	resetUsers()

	// Setup
	invalidJSON := []byte(`{"fullName": invalid}`)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	// Execute
	UpdateUser(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["error"])
}
