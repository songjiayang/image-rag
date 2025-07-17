package handlers

import (
	"net/http"
	"time"

	"image-rag-backend/internal/milvus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db     *gorm.DB
	milvus *milvus.Client
}

func NewHealthHandler(db *gorm.DB, milvus *milvus.Client) *HealthHandler {
	return &HealthHandler{
		db:     db,
		milvus: milvus,
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		response.Services["database"] = "unhealthy"
		response.Status = "degraded"
	} else if err := sqlDB.Ping(); err != nil {
		response.Services["database"] = "unhealthy"
		response.Status = "degraded"
	} else {
		response.Services["database"] = "healthy"
	}

	// Check Milvus connection
	if err := h.milvus.Ping(); err != nil {
		response.Services["milvus"] = "unhealthy"
		response.Status = "degraded"
	} else {
		response.Services["milvus"] = "healthy"
	}

	// Check Doubao API (basic connectivity)
	response.Services["doubao"] = "healthy" // We'll implement actual check later

	if response.Status == "healthy" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// More comprehensive health check for Kubernetes readiness
	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	// Check all critical services
	sqlDB, dbErr := h.db.DB()
	if dbErr == nil {
		dbErr = sqlDB.Ping()
	}
	milvusErr := h.milvus.Ping()

	if dbErr == nil {
		response.Services["database"] = "ready"
	} else {
		response.Services["database"] = "not ready"
		response.Status = "not ready"
	}

	if milvusErr == nil {
		response.Services["milvus"] = "ready"
	} else {
		response.Services["milvus"] = "not ready"
		response.Status = "not ready"
	}

	if response.Status == "ready" {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusServiceUnavailable, response)
	}
}

func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// Simple liveness check - just returns 200 if service is running
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
