package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"manhwa-tools-backend/internal/api/tools"
)

// SetupRouter initializes the Gin router and registers all routes
func SetupRouter(modelPath string) *gin.Engine {
	r := gin.Default()
	
	// Middleware
	r.Use(cors.Default())

	// Static Frontend
	r.Static("/", "../frontend")

	// API Group
	apiGroup := r.Group("/api")
	
	// Register Tools
	tools.RegisterEraserTool(apiGroup, modelPath)
	tools.RegisterTranslatorTool(apiGroup, modelPath)
	
	// In the future, just add:
	// tools.RegisterUpscalerTool(apiGroup)

	return r
}
