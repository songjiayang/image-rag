package services

import (
	"time"

	"image-rag-backend/internal/models"

	"gorm.io/gorm"
)

type StatsService struct {
	db *gorm.DB
}

func NewStatsService(db *gorm.DB) *StatsService {
	return &StatsService{db: db}
}

func (s *StatsService) GetDashboardStats() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{}

	// Get total records count
	var totalRecords int64
	if err := s.db.Model(&models.Record{}).Count(&totalRecords).Error; err != nil {
		return nil, err
	}
	stats.TotalRecords = totalRecords

	// Get total images count
	var totalImages int64
	if err := s.db.Model(&models.Image{}).Count(&totalImages).Error; err != nil {
		return nil, err
	}
	stats.TotalImages = totalImages

	// Get today's date boundaries
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Get today's records count
	var todayRecords int64
	if err := s.db.Model(&models.Record{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&todayRecords).Error; err != nil {
		return nil, err
	}
	stats.TodayRecords = todayRecords

	// Get today's images count
	var todayImages int64
	if err := s.db.Model(&models.Image{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Count(&todayImages).Error; err != nil {
		return nil, err
	}
	stats.TodayImages = todayImages

	return stats, nil
}
