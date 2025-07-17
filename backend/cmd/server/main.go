package main

import (
	"context"
	"image-rag-backend/internal/api"
	"image-rag-backend/internal/config"
	"image-rag-backend/internal/database"
	"image-rag-backend/internal/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Ensure upload directories exist
	if err := services.EnsureDirectoryExists(cfg.Upload.Path); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}
	if err := services.EnsureDirectoryExists(filepath.Join(cfg.Upload.Path, "temp")); err != nil {
		log.Fatal("Failed to create temp directory:", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") != "release" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Setup routes
	api.SetupRoutes(router, cfg)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
