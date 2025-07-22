package handlers

import (
	"net/http"

	"image-rag-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsService *services.StatsService
}

func NewStatsHandler(statsService *services.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetDashboardStats returns dashboard statistics
// @Summary Get dashboard statistics
// @Description Get comprehensive dashboard statistics including total records, images, and today's counts
// @Tags Stats
// @Produce json
// @Success 200 {object} map[string]interface{} "{"data": {"total_records": 100, "total_images": 250, "today_records": 5, "today_images": 12}}"
// @Failure 500 {object} map[string]string "{"error": "Failed to fetch dashboard statistics"}"
// @Router /api/v1/stats [get]
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch dashboard statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}
