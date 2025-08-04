package main

import (
	"log"
	"tokobiru/config"
	"tokobiru/database"
	"tokobiru/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Connect to MongoDB
	dbClient := database.ConnectDB(cfg.MongoURI, cfg.MongoDatabase)

	// Set Gin to release mode for production
	// gin.SetMode(gin.ReleaseMode)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, dbClient)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
