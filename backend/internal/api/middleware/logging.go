package middleware

import (
	"fmt"
	"time"

	"image-rag-backend/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Start timer
		start := time.Now()

		// Log request
		log.Info("[%s] %s %s - Started", requestID, c.Request.Method, c.Request.URL.Path)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log response
		status := c.Writer.Status()
		if status >= 400 {
			log.ErrorWithContext(map[string]interface{}{
				"request_id": requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     status,
				"latency":    latency.String(),
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
				"error":      c.Errors.String(),
			}, "HTTP request failed")
		} else {
			log.Info("[%s] %s %s - %d (%s)", requestID, c.Request.Method, c.Request.URL.Path, status, latency.String())
		}
	}
}

func ErrorHandlerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID, _ := c.Get("request_id")

		// Log the panic
		log.ErrorWithContext(map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"panic":      fmt.Sprintf("%v", recovered),
			"stack":      string(getStackTrace()),
		}, "Application panic recovered")

		// Return error response
		c.JSON(500, gin.H{
			"error":      "Internal server error",
			"request_id": requestID,
			"message":    "Something went wrong, please try again later",
		})
		c.Abort()
	})
}

func getStackTrace() []byte {
	// In a real implementation, you'd capture the actual stack trace
	// For now, return a placeholder
	return []byte("stack trace not implemented")
}
