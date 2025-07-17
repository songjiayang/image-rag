package models

import (
	"time"
)

type Record struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Description string    `json:"description" gorm:"type:text"`
	Images      []Image   `json:"images" gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Image struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RecordID  uint      `json:"record_id" gorm:"not null;index"`
	Filename  string    `json:"filename" gorm:"not null;size:255"`
	Path      string    `json:"path" gorm:"not null;size:500"`
	VectorID  string    `json:"vector_id" gorm:"not null;size:100;index"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRecordRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateRecordRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RecordResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Images      []ImageResponse `json:"images"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ImageResponse struct {
	ID       uint   `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
}
