package api

import (
	_ "image-rag-backend/docs"
	"image-rag-backend/internal/api/handlers"
	"image-rag-backend/internal/api/middleware"
	"image-rag-backend/internal/config"
	"image-rag-backend/internal/database"
	"image-rag-backend/internal/milvus"
	"image-rag-backend/internal/services"
	"log"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	// Initialize services
	recordService := services.NewRecordService()
	vectorService, err := services.NewVectorService(cfg)
	if err != nil {
		log.Fatal("Failed to initialize vector service:", err)
	}

	// Initialize handlers
	recordHandler := handlers.NewRecordHandler(recordService, vectorService)
	searchHandler := handlers.NewSearchHandler(recordService, vectorService)

	// Global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.DefaultRateLimit())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Set max multipart memory to 64MB
	router.MaxMultipartMemory = 64 << 20 // 64 MB

	// API routes
	api := router.Group("/api/v1")

	// Swagger documentation
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Initialize database and milvus clients for health checks
	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	milvusClient, err := milvus.NewClient(&cfg.Milvus)
	if err != nil {
		log.Fatal("Failed to connect to milvus:", err)
	}

	healthHandler := handlers.NewHealthHandler(database.DB, milvusClient)

	// Health check endpoints
	api.GET("/health", healthHandler.HealthCheck)
	api.GET("/health/ready", healthHandler.ReadinessCheck)
	api.GET("/health/live", healthHandler.LivenessCheck)

	// Records routes
	api.POST("/records", recordHandler.CreateRecord)
	api.GET("/records", recordHandler.GetRecords)
	api.GET("/records/:id", recordHandler.GetRecord)
	api.PUT("/records/:id", recordHandler.UpdateRecord)
	api.DELETE("/records/:id", recordHandler.DeleteRecord)

	// Image management routes
	api.POST("/records/:id/images", recordHandler.AddImageToRecord)
	api.DELETE("/images/:image_id", recordHandler.DeleteImage)

	// Search routes
	api.POST("/search", searchHandler.SearchImages)
	api.GET("/search/similar/:id", searchHandler.FindSimilar)
	api.POST("/search/advanced", searchHandler.AdvancedSearch)
	api.GET("/search/by-vector/:vector_id", searchHandler.GetImageByVectorID)

	// Serve uploaded images
	router.Static("/uploads", "./uploads")
}

// Note: The services.NewConfig() should be properly initialized from main.go
// This is a placeholder for the route setup
