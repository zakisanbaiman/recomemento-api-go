package main

import (
	"log"
	"os"

	"recomemento-api-go/database"
	_ "recomemento-api-go/docs" // Swagger docs
	"recomemento-api-go/handlers"
	"recomemento-api-go/models"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Recomemento API
// @version 1.0
// @description 本の推薦システムのバックエンドAPI
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3001
// @BasePath /
func main() {
	// Database initialization
	dbPath := "./data/books.db"
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		if dbURL == "file:./prisma/dev.db" {
			dbPath = "./prisma/dev.db"
		} else {
			dbPath = dbURL
		}
	}

	db, err := database.InitDatabase(dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Seed database
	if err := database.SeedDatabase(db); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	// Initialize repositories
	bookRepo := models.NewBookRepository(db)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookRepo)

	// Initialize Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Recomemento API is running"})
	})

	// API routes
	api := r.Group("/")
	{
		// Book routes
		api.POST("/books", bookHandler.CreateBook)
		api.GET("/books", bookHandler.GetAllBooks)
		api.GET("/books/:id", bookHandler.GetBookByID)
		api.PATCH("/books/:id", bookHandler.UpdateBook)
		api.DELETE("/books/:id", bookHandler.DeleteBook)
		api.POST("/books/recommend", bookHandler.RecommendBook)
	}

	// Swagger documentation
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// OpenAPI JSON endpoint (compatible with original API)
	r.GET("/api-json", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.File("./docs/swagger.json")
	})

	// Get port from environment or default to 3001
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Swagger UI available at: http://localhost:%s/api-docs/", port)
	log.Printf("Health check available at: http://localhost:%s/health", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 