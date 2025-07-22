package models

type DashboardStats struct {
	TotalRecords int64 `json:"total_records"`
	TotalImages  int64 `json:"total_images"`
	TodayRecords int64 `json:"today_records"`
	TodayImages  int64 `json:"today_images"`
}
