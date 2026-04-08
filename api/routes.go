package api

import (
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"userprofile-api/controllers"
)

// SetupRouter configures the API routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Get the absolute path to the templates directory
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(b))
	templatesPath := filepath.Join(basePath, "templates/*")

	// Setup template rendering
	router.LoadHTMLGlob(templatesPath)

	uc := controllers.NewUserController(db)

	// Root handler shows a nice HTML table of all users
	router.GET("/", uc.HomePageHandler)

	// API version group
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", uc.GetUsers)
			users.GET("/:id", uc.GetUser)
			users.POST("", uc.CreateUser)
			users.PUT("/:id", uc.UpdateUser)
			users.DELETE("/:id", uc.DeleteUser)
		}
	}

	return router
}
